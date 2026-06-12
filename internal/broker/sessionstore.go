package broker

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/storage"
	"github.com/mochi-mqtt/server/v2/packets"
	"go.etcd.io/bbolt"
)

// bbolt bucket names.
var (
	bucketClients       = []byte("clients")
	bucketSubscriptions = []byte("subscriptions")
	bucketInflight      = []byte("inflight")
	bucketRetained      = []byte("retained")
)

// keySep joins composite keys ("<clientID>\x00<suffix>"). NUL cannot appear
// in MQTT client identifiers or topic filters, so prefix scans are exact.
const keySep = byte(0)

// bboltOpenTimeout bounds the file-lock wait when opening the session store,
// so a second broker instance on the same file fails fast instead of
// hanging.
const bboltOpenTimeout = 2 * time.Second

// storedMessage is the persisted form of an inflight or retained message.
// It extends mochi's storage.Message with the fields mochi's restore path
// loses but Astrate needs:
//
//   - Expiry/ProtocolVersion: per-message expiry on the offline queue
//     (docs/DESIGN.md §2.5, §3.4) — expired messages are dropped at load
//     time, since storage.Message cannot carry them back into the server;
//   - Seq: a monotonic insertion sequence. mochi orders inflight replay by
//     uint16(Created) (whole seconds), so messages queued within the same
//     second would replay in arbitrary order after a restart; Seq lets the
//     load path normalize Created into a strictly-increasing sequence that
//     preserves true queueing order (the ordered-replay guarantee of
//     docs/ROADMAP.md §6).
type storedMessage struct {
	storage.Message
	Expiry          int64  `json:"expiry,omitempty"`
	ProtocolVersion byte   `json:"pv,omitempty"`
	Seq             uint64 `json:"seq,omitempty"`
}

// sessionStore is the bbolt-backed mochi storage hook (docs/ROADMAP.md §6
// file 5.5): client sessions, subscriptions, the QoS >= 1 inflight/offline
// queue, and retained messages survive broker restarts in a single file.
type sessionStore struct {
	mqtt.HookBase
	path string
	db   *bbolt.DB
	srv  *mqtt.Server // for pre-inherit inflight pruning; set via attach
	log  *slog.Logger
	now  func() time.Time
}

func newSessionStore(path string, log *slog.Logger) *sessionStore {
	return &sessionStore{path: path, log: log, now: time.Now}
}

// attach hands the store its server handle (used by OnSessionEstablish).
// Called by Broker.New after mqtt.New, before Serve.
func (s *sessionStore) attach(srv *mqtt.Server) { s.srv = srv }

// ID implements mqtt.Hook.
func (s *sessionStore) ID() string { return "astrate-session-store" }

// Provides implements mqtt.Hook.
func (s *sessionStore) Provides(b byte) bool {
	switch b {
	case mqtt.OnSessionEstablish, mqtt.OnSessionEstablished, mqtt.OnDisconnect,
		mqtt.OnSubscribed, mqtt.OnUnsubscribed,
		mqtt.OnRetainMessage, mqtt.OnRetainedExpired,
		mqtt.OnQosPublish, mqtt.OnQosComplete, mqtt.OnQosDropped,
		mqtt.OnClientExpired,
		mqtt.StoredClients, mqtt.StoredSubscriptions, mqtt.StoredInflightMessages,
		mqtt.StoredRetainedMessages, mqtt.StoredSysInfo:
		return true
	}
	return false
}

// Init implements mqtt.Hook: it opens the bbolt file and creates buckets.
func (s *sessionStore) Init(config any) error {
	if config != nil {
		return mqtt.ErrInvalidConfigType
	}
	db, err := bbolt.Open(s.path, 0o600, &bbolt.Options{Timeout: bboltOpenTimeout})
	if err != nil {
		return fmt.Errorf("broker: opening session store %q: %w", s.path, err)
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		for _, b := range [][]byte{bucketClients, bucketSubscriptions, bucketInflight, bucketRetained} {
			if _, err := tx.CreateBucketIfNotExists(b); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		_ = db.Close()
		return fmt.Errorf("broker: preparing session store buckets: %w", err)
	}
	s.db = db
	return nil
}

// Stop implements mqtt.Hook: it closes (and releases the file lock on) the
// bbolt database. Called from mqtt.Server.Close.
func (s *sessionStore) Stop() error {
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

// --- write-side hooks ---

// OnSessionEstablished persists the client record.
func (s *sessionStore) OnSessionEstablished(cl *mqtt.Client, _ packets.Packet) {
	props := cl.Properties.Props.Copy(false)
	rec := storage.Client{
		ID:              cl.ID,
		T:               storage.ClientKey,
		Remote:          cl.Net.Remote,
		Listener:        cl.Net.Listener,
		Username:        cl.Properties.Username,
		Clean:           cl.Properties.Clean,
		ProtocolVersion: cl.Properties.ProtocolVersion,
		Properties: storage.ClientProperties{
			SessionExpiryInterval:     props.SessionExpiryInterval,
			SessionExpiryIntervalFlag: props.SessionExpiryIntervalFlag,
			ReceiveMaximum:            props.ReceiveMaximum,
			TopicAliasMaximum:         props.TopicAliasMaximum,
			MaximumPacketSize:         props.MaximumPacketSize,
		},
	}
	s.put(bucketClients, []byte(cl.ID), rec)
}

// OnDisconnect drops expiring sessions; persistent sessions stay on disk.
func (s *sessionStore) OnDisconnect(cl *mqtt.Client, _ error, expire bool) {
	if !expire || errors.Is(cl.StopCause(), packets.ErrSessionTakenOver) {
		return
	}
	s.dropClient(cl.ID)
}

// OnClientExpired drops a session that outlived its expiry interval.
func (s *sessionStore) OnClientExpired(cl *mqtt.Client) {
	s.dropClient(cl.ID)
}

// OnSubscribed persists granted subscriptions (denied filters carry a
// failure reason code and are skipped).
func (s *sessionStore) OnSubscribed(cl *mqtt.Client, pk packets.Packet, reasonCodes []byte) {
	for i := 0; i < len(pk.Filters) && i < len(reasonCodes); i++ {
		if reasonCodes[i] >= packets.ErrUnspecifiedError.Code {
			continue
		}
		f := pk.Filters[i]
		rec := storage.Subscription{
			T:                 storage.SubscriptionKey,
			Client:            cl.ID,
			Filter:            f.Filter,
			Identifier:        f.Identifier,
			Qos:               reasonCodes[i],
			RetainHandling:    f.RetainHandling,
			RetainAsPublished: f.RetainAsPublished,
			NoLocal:           f.NoLocal,
		}
		s.put(bucketSubscriptions, compositeKey(cl.ID, []byte(f.Filter)), rec)
	}
}

// OnUnsubscribed removes subscriptions.
func (s *sessionStore) OnUnsubscribed(cl *mqtt.Client, pk packets.Packet) {
	for _, f := range pk.Filters {
		s.delete(bucketSubscriptions, compositeKey(cl.ID, []byte(f.Filter)))
	}
}

// OnRetainMessage tracks the retained message set (r == -1 clears).
func (s *sessionStore) OnRetainMessage(_ *mqtt.Client, pk packets.Packet, r int64) {
	if r == -1 {
		s.delete(bucketRetained, []byte(pk.TopicName))
		return
	}
	s.put(bucketRetained, []byte(pk.TopicName), messageRecord(pk, "", storage.RetainedKey, 0))
}

// OnRetainedExpired removes an expired retained message.
func (s *sessionStore) OnRetainedExpired(filter string) {
	s.delete(bucketRetained, []byte(filter))
}

// OnQosPublish upserts an inflight (or offline-queued) message. Resends
// re-fire this hook; the original Seq (queueing order) is preserved.
func (s *sessionStore) OnQosPublish(cl *mqtt.Client, pk packets.Packet, _ int64, _ int) {
	if s.db == nil {
		return
	}
	key := inflightKey(cl.ID, pk.PacketID)
	err := s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketInflight)
		seq := uint64(0)
		if prev := b.Get(key); prev != nil {
			var old storedMessage
			if json.Unmarshal(prev, &old) == nil {
				seq = old.Seq
			}
		}
		if seq == 0 {
			var err error
			if seq, err = b.NextSequence(); err != nil {
				return err
			}
		}
		buf, err := json.Marshal(messageRecord(pk, cl.ID, storage.InflightKey, seq))
		if err != nil {
			return err
		}
		return b.Put(key, buf)
	})
	if err != nil {
		s.log.Error("session store: persisting inflight message", "client", cl.ID, "error", err)
	}
}

// OnQosComplete removes a completed inflight message.
func (s *sessionStore) OnQosComplete(cl *mqtt.Client, pk packets.Packet) {
	s.delete(bucketInflight, inflightKey(cl.ID, pk.PacketID))
}

// OnQosDropped removes a dropped (e.g. expired) inflight message.
func (s *sessionStore) OnQosDropped(cl *mqtt.Client, pk packets.Packet) {
	s.delete(bucketInflight, inflightKey(cl.ID, pk.PacketID))
}

// OnSessionEstablish runs before an existing session is inherited: it prunes
// queued messages whose per-message expiry elapsed while this broker
// instance was running with the session offline. (mochi's own sweeper only
// sees expiry on packets that carry it in memory; messages restored from
// disk no longer do — see storedMessage.)
func (s *sessionStore) OnSessionEstablish(cl *mqtt.Client, _ packets.Packet) {
	if s.db == nil || s.srv == nil {
		return
	}
	old, ok := s.srv.Clients.Get(cl.ID)
	if !ok || old.State.Inflight.Len() == 0 {
		return
	}
	now := s.now().Unix()
	prefix := compositeKey(cl.ID, nil)
	err := s.db.Update(func(tx *bbolt.Tx) error {
		c := tx.Bucket(bucketInflight).Cursor()
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			var rec storedMessage
			if json.Unmarshal(v, &rec) != nil {
				continue
			}
			if rec.Expiry > 0 && rec.Expiry <= now {
				old.State.Inflight.Delete(rec.PacketID)
				if err := c.Delete(); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		s.log.Error("session store: pruning expired inflight messages", "client", cl.ID, "error", err)
	}
}

// --- restore-side hooks ---

// StoredClients implements mqtt.Hook.
func (s *sessionStore) StoredClients() ([]storage.Client, error) {
	var out []storage.Client
	err := s.view(bucketClients, func(_, v []byte) error {
		var rec storage.Client
		if err := json.Unmarshal(v, &rec); err != nil {
			return err
		}
		out = append(out, rec)
		return nil
	})
	return out, err
}

// StoredSubscriptions implements mqtt.Hook.
func (s *sessionStore) StoredSubscriptions() ([]storage.Subscription, error) {
	var out []storage.Subscription
	err := s.view(bucketSubscriptions, func(_, v []byte) error {
		var rec storage.Subscription
		if err := json.Unmarshal(v, &rec); err != nil {
			return err
		}
		out = append(out, rec)
		return nil
	})
	return out, err
}

// StoredInflightMessages implements mqtt.Hook: it returns the surviving
// offline queue, dropping (from disk too) messages whose expiry elapsed
// while the broker was down, and normalizing Created per client so replay
// order matches queueing order (see storedMessage).
func (s *sessionStore) StoredInflightMessages() ([]storage.Message, error) {
	if s.db == nil {
		return nil, storage.ErrDBFileNotOpen
	}
	now := s.now().Unix()
	var live []storedMessage
	err := s.db.Update(func(tx *bbolt.Tx) error {
		c := tx.Bucket(bucketInflight).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var rec storedMessage
			if err := json.Unmarshal(v, &rec); err != nil {
				return err
			}
			if rec.Expiry > 0 && rec.Expiry <= now {
				if err := c.Delete(); err != nil {
					return err
				}
				continue
			}
			live = append(live, rec)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return normalizeInflightOrder(live), nil
}

// StoredRetainedMessages implements mqtt.Hook, dropping expired retained
// messages at load time.
func (s *sessionStore) StoredRetainedMessages() ([]storage.Message, error) {
	if s.db == nil {
		return nil, storage.ErrDBFileNotOpen
	}
	now := s.now().Unix()
	var out []storage.Message
	err := s.db.Update(func(tx *bbolt.Tx) error {
		c := tx.Bucket(bucketRetained).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var rec storedMessage
			if err := json.Unmarshal(v, &rec); err != nil {
				return err
			}
			if rec.Expiry > 0 && rec.Expiry <= now {
				if err := c.Delete(); err != nil {
					return err
				}
				continue
			}
			out = append(out, rec.Message)
		}
		return nil
	})
	return out, err
}

// StoredSysInfo implements mqtt.Hook. Astrate does not persist $SYS
// counters; observability is Prometheus' job (docs/DESIGN.md §5.2).
func (s *sessionStore) StoredSysInfo() (storage.SystemInfo, error) {
	return storage.SystemInfo{}, nil
}

// --- helpers ---

// normalizeInflightOrder sorts each client's queue by Seq and rewrites
// Created into a strictly-increasing sequence (clamped forward only), so
// mochi's uint16(Created) replay sort reproduces queueing order. Records are
// assumed pre-sorted by key, i.e. grouped by client; insertion sort by Seq
// keeps it dependency-free.
func normalizeInflightOrder(recs []storedMessage) []storage.Message {
	byClient := map[string][]storedMessage{}
	var clients []string
	for _, r := range recs {
		if _, ok := byClient[r.Client]; !ok {
			clients = append(clients, r.Client)
		}
		byClient[r.Client] = append(byClient[r.Client], r)
	}

	out := make([]storage.Message, 0, len(recs))
	for _, id := range clients {
		group := byClient[id]
		for i := 1; i < len(group); i++ { // insertion sort by Seq
			for j := i; j > 0 && group[j-1].Seq > group[j].Seq; j-- {
				group[j-1], group[j] = group[j], group[j-1]
			}
		}
		var prev int64
		for i, r := range group {
			msg := r.Message
			if i > 0 && msg.Created <= prev {
				msg.Created = prev + 1
			}
			prev = msg.Created
			out = append(out, msg)
		}
	}
	return out
}

// messageRecord converts a packet to its persisted form.
func messageRecord(pk packets.Packet, clientID, typ string, seq uint64) storedMessage {
	return storedMessage{
		Message: storage.Message{
			T:         typ,
			Client:    clientID,
			Origin:    pk.Origin,
			TopicName: pk.TopicName,
			Payload:   pk.Payload,
			Created:   pk.Created,
			PacketID:  pk.PacketID,
			FixedHeader: packets.FixedHeader{
				Type:   pk.FixedHeader.Type,
				Qos:    pk.FixedHeader.Qos,
				Retain: pk.FixedHeader.Retain,
				Dup:    pk.FixedHeader.Dup,
			},
			Properties: storage.MessageProperties{
				PayloadFormat:         pk.Properties.PayloadFormat,
				PayloadFormatFlag:     pk.Properties.PayloadFormatFlag,
				MessageExpiryInterval: pk.Properties.MessageExpiryInterval,
				ContentType:           pk.Properties.ContentType,
				ResponseTopic:         pk.Properties.ResponseTopic,
				CorrelationData:       pk.Properties.CorrelationData,
				User:                  pk.Properties.User,
			},
		},
		Expiry:          pk.Expiry,
		ProtocolVersion: pk.ProtocolVersion,
		Seq:             seq,
	}
}

// compositeKey builds "<clientID>\x00<suffix>".
func compositeKey(clientID string, suffix []byte) []byte {
	key := make([]byte, 0, len(clientID)+1+len(suffix))
	key = append(key, clientID...)
	key = append(key, keySep)
	return append(key, suffix...)
}

// inflightKey builds the per-message inflight key.
func inflightKey(clientID string, packetID uint16) []byte {
	var pid [2]byte
	binary.BigEndian.PutUint16(pid[:], packetID)
	return compositeKey(clientID, pid[:])
}

// put marshals and stores one record.
func (s *sessionStore) put(bucket, key []byte, rec any) {
	if s.db == nil {
		return
	}
	buf, err := json.Marshal(rec)
	if err != nil {
		s.log.Error("session store: encoding record", "bucket", string(bucket), "error", err)
		return
	}
	if err := s.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucket).Put(key, buf)
	}); err != nil {
		s.log.Error("session store: writing record", "bucket", string(bucket), "error", err)
	}
}

// delete removes one record.
func (s *sessionStore) delete(bucket, key []byte) {
	if s.db == nil {
		return
	}
	if err := s.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucket).Delete(key)
	}); err != nil {
		s.log.Error("session store: deleting record", "bucket", string(bucket), "error", err)
	}
}

// dropClient removes a client and all its dependent records.
func (s *sessionStore) dropClient(clientID string) {
	if s.db == nil {
		return
	}
	prefix := compositeKey(clientID, nil)
	err := s.db.Update(func(tx *bbolt.Tx) error {
		if err := tx.Bucket(bucketClients).Delete([]byte(clientID)); err != nil {
			return err
		}
		for _, bucket := range [][]byte{bucketSubscriptions, bucketInflight} {
			c := tx.Bucket(bucket).Cursor()
			for k, _ := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = c.Next() {
				if err := c.Delete(); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		s.log.Error("session store: dropping client", "client", clientID, "error", err)
	}
}

// view iterates a bucket read-only.
func (s *sessionStore) view(bucket []byte, fn func(k, v []byte) error) error {
	if s.db == nil {
		return storage.ErrDBFileNotOpen
	}
	return s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucket).ForEach(fn)
	})
}
