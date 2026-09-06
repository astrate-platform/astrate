package housekeeping

import (
	"context"
	"errors"
	"testing"

	"github.com/astrate-platform/astrate/internal/store"
)

// fakeStore implements the realmStore port the service needs, recording every
// call so tests can assert what the service did (and did not) touch.
type fakeStore struct {
	realm    *store.Realm
	realmErr error

	total, connected int64
	statsErr         error

	createErr error
	updateErr error
	deleteErr error

	created []store.NewRealm
	deleted []string
	got     []string // names looked up via GetRealmByName
}

var _ realmStore = (*fakeStore)(nil)

func (f *fakeStore) CreateRealm(_ context.Context, nr store.NewRealm) (*store.Realm, error) {
	f.created = append(f.created, nr)
	if f.createErr != nil {
		return nil, f.createErr
	}
	return &store.Realm{
		Name:                              nr.Name,
		JWTPublicKeysPEM:                  nr.JWTPublicKeysPEM,
		DeviceRegistrationLimit:           nr.DeviceRegistrationLimit,
		DatastreamMaximumStorageRetention: nr.DatastreamMaximumStorageRetention,
	}, nil
}

func (f *fakeStore) GetRealmByName(_ context.Context, name string) (*store.Realm, error) {
	f.got = append(f.got, name)
	if f.realmErr != nil {
		return nil, f.realmErr
	}
	if f.realm == nil {
		return nil, store.ErrNotFound
	}
	return f.realm, nil
}

func (f *fakeStore) UpdateRealm(_ context.Context, _ string, _ store.RealmPatch) error {
	return f.updateErr
}

func (f *fakeStore) ListRealms(_ context.Context) ([]store.Realm, error) {
	return nil, nil
}

func (f *fakeStore) DeviceStats(_ context.Context, _ int16) (total, connected int64, err error) {
	if f.statsErr != nil {
		return 0, 0, f.statsErr
	}
	return f.total, f.connected, nil
}

func (f *fakeStore) DeleteRealm(_ context.Context, name string) error {
	f.deleted = append(f.deleted, name)
	return f.deleteErr
}

// fakeSealer implements the sealer port, recording every sealed input.
type fakeSealer struct {
	sealed [][]byte
	err    error
}

var _ sealer = (*fakeSealer)(nil)

func (f *fakeSealer) Seal(plaintext []byte) ([]byte, error) {
	f.sealed = append(f.sealed, plaintext)
	if f.err != nil {
		return nil, f.err
	}
	return append([]byte(nil), plaintext...), nil
}

// TestCreateRealmValidation covers the ErrValidation rule of CreateRealm: a
// blank name, a blank JWT key, a negative registration limit, and a negative
// retention all fail before the store or sealer is touched.
func TestCreateRealmValidation(t *testing.T) {
	ctx := context.Background()
	negLimit := int32(-1)
	negRetention := int64(-1)

	for _, tc := range []struct {
		name, realmName, jwtKey string
		regLimit                *int32
		retention               *int64
	}{
		{"blank name", "", "key", nil, nil},
		{"blank jwt key", "plant", "", nil, nil},
		{"negative registration limit", "plant", "key", &negLimit, nil},
		{"negative retention", "plant", "key", nil, &negRetention},
	} {
		t.Run(tc.name, func(t *testing.T) {
			st := &fakeStore{}
			sl := &fakeSealer{}
			svc := NewService(st, sl, nil, nil)

			if _, err := svc.CreateRealm(ctx, tc.realmName, tc.jwtKey, tc.regLimit, tc.retention); !errors.Is(err, ErrValidation) {
				t.Fatalf("CreateRealm(%q, %q) error = %v, want ErrValidation", tc.realmName, tc.jwtKey, err)
			}
			if len(st.created) != 0 {
				t.Errorf("store.CreateRealm called %d times on a rejected create", len(st.created))
			}
			if len(sl.sealed) != 0 {
				t.Errorf("sealer.Seal called %d times on a rejected create", len(sl.sealed))
			}
		})
	}
}

// TestCreateRealmDefaultRetention covers #73 at the service level: the
// configured default is injected only when the caller omits retention; an
// explicit value always wins.
func TestCreateRealmDefaultRetention(t *testing.T) {
	ctx := context.Background()
	def := int64(3600)

	t.Run("omitted injects default", func(t *testing.T) {
		st := &fakeStore{}
		svc := NewService(st, &fakeSealer{}, nil, nil).
			WithDefaultDatastreamMaximumStorageRetention(&def)

		rv, err := svc.CreateRealm(ctx, "pico", "key", nil, nil)
		if err != nil {
			t.Fatalf("CreateRealm: %v", err)
		}
		if len(st.created) != 1 {
			t.Fatalf("store.CreateRealm calls = %d, want 1", len(st.created))
		}
		if got := st.created[0].DatastreamMaximumStorageRetention; got == nil || *got != def {
			t.Errorf("store received retention = %v, want %d", got, def)
		}
		if got := rv.DatastreamMaximumStorageRetention; got == nil || *got != def {
			t.Errorf("view retention = %v, want %d", got, def)
		}
	})

	t.Run("explicit value beats default", func(t *testing.T) {
		st := &fakeStore{}
		explicit := int64(60)
		svc := NewService(st, &fakeSealer{}, nil, nil).
			WithDefaultDatastreamMaximumStorageRetention(&def)

		if _, err := svc.CreateRealm(ctx, "micro", "key", nil, &explicit); err != nil {
			t.Fatalf("CreateRealm: %v", err)
		}
		if got := st.created[0].DatastreamMaximumStorageRetention; got == nil || *got != explicit {
			t.Errorf("store received retention = %v, want %d", got, explicit)
		}
	})

	t.Run("no default stays nil", func(t *testing.T) {
		st := &fakeStore{}
		svc := NewService(st, &fakeSealer{}, nil, nil)

		if _, err := svc.CreateRealm(ctx, "nano", "key", nil, nil); err != nil {
			t.Fatalf("CreateRealm: %v", err)
		}
		if st.created[0].DatastreamMaximumStorageRetention != nil {
			t.Errorf("store received retention %v, want nil", *st.created[0].DatastreamMaximumStorageRetention)
		}
	})
}

// TestDeleteRealmGating covers #75 at the service level: the deletion-disabled
// flag answers ErrDeletionDisabled before any store access, and a realm with
// connected devices answers ErrConnectedDevicesPresent without deleting.
func TestDeleteRealmGating(t *testing.T) {
	ctx := context.Background()

	t.Run("deletion disabled short-circuits", func(t *testing.T) {
		st := &fakeStore{}
		svc := NewService(st, &fakeSealer{}, nil, nil).WithRealmDeletionDisabled(true)

		// The gate precedes the existence check: the unknown realm would
		// otherwise surface store.ErrNotFound.
		if err := svc.DeleteRealm(ctx, "whatever"); !errors.Is(err, ErrDeletionDisabled) {
			t.Fatalf("DeleteRealm = %v, want ErrDeletionDisabled", err)
		}
		if len(st.got) != 0 {
			t.Errorf("gate must precede the store lookup, got %d lookups", len(st.got))
		}
	})

	t.Run("connected device blocks delete", func(t *testing.T) {
		st := &fakeStore{realm: &store.Realm{ID: 7, Name: "busy"}, connected: 3}
		svc := NewService(st, &fakeSealer{}, nil, nil)

		if err := svc.DeleteRealm(ctx, "busy"); !errors.Is(err, ErrConnectedDevicesPresent) {
			t.Fatalf("DeleteRealm = %v, want ErrConnectedDevicesPresent", err)
		}
		if len(st.deleted) != 0 {
			t.Errorf("store.DeleteRealm called on a gated delete")
		}
	})

	t.Run("no connected devices deletes", func(t *testing.T) {
		st := &fakeStore{realm: &store.Realm{ID: 7, Name: "calm"}}
		svc := NewService(st, &fakeSealer{}, nil, nil)

		if err := svc.DeleteRealm(ctx, "calm"); err != nil {
			t.Fatalf("DeleteRealm = %v, want nil", err)
		}
		if len(st.deleted) != 1 || st.deleted[0] != "calm" {
			t.Errorf("store.DeleteRealm calls = %v, want [calm]", st.deleted)
		}
	})
}
