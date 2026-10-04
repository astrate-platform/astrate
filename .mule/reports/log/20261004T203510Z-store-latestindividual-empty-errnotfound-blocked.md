slug: store-latestindividual-empty-errnotfound
verdict: blocked
at:  945371f
ran: 2026-10-04T20:35:10Z on DietPi in 1201s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md

$ grep -rn "LatestIndividual" /root/astrate-mule/internal/store
/root/astrate-mule/internal/store/datastreams_test.go:175:	t.Run("LatestIndividual", func(t *testing.T) {
/root/astrate-mule/internal/store/datastreams_test.go:196:		got, err := s.LatestIndividual(ctx, realm.ID, device, si.ID, "/d")
/root/astrate-mule/internal/store/datastreams_test.go:198:			t.Fatalf("LatestIndividual: %v", err)
/root/astrate-mule/internal/store/datastreams_test.go:208:		_, err = s.LatestIndividual(ctx, realm.ID, device, si.ID, "/missing")
/root/astrate-mule/internal/store/datastreams.go:243:// LatestIndividual reads the newest sample of one individual-datastream
/root/astrate-mule/internal/store/datastreams.go:247:func (s *Store) LatestIndividual(ctx context.Context, realmID int16, deviceID deviceid.ID, interfaceID int64, path string) (*IndividualRow, error) {

→ Read internal/store/datastreams.go
