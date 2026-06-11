package payload

import (
	"errors"
	"fmt"
	"math"
	"testing"
	"time"
)

// TestDoubleFromInt64 pins the "int64 → double lossless only" rule: 2^53 is
// the last contiguous integer float64 represents exactly.
func TestDoubleFromInt64(t *testing.T) {
	cases := []struct {
		in   int64
		want float64
		ok   bool
	}{
		{0, 0, true},
		{42, 42, true},
		{-42, -42, true},
		{1 << 53, 9007199254740992, true},           // 2^53: exactly representable
		{1<<53 + 1, 0, false},                       // 2^53+1: rounds to 2^53 → lossy
		{-(1 << 53), -9007199254740992, true},       // -2^53
		{-(1<<53 + 1), 0, false},                    // -(2^53+1)
		{1<<53 + 2, 9007199254740994, true},         // even values above 2^53 stay exact
		{math.MaxInt64, 0, false},                   // 2^63-1 rounds to 2^63
		{math.MinInt64, -9223372036854775808, true}, // -2^63 is a power of two: exact
	}
	for _, tc := range cases {
		got, err := doubleFromInt64(tc.in)
		if tc.ok && (err != nil || got != tc.want) {
			t.Errorf("doubleFromInt64(%d) = %v, %v; want %v, nil", tc.in, got, err, tc.want)
		}
		if !tc.ok {
			if err == nil {
				t.Errorf("doubleFromInt64(%d) accepted lossy widening", tc.in)
			} else if ReasonOf(err) != ReasonTypeMismatch {
				t.Errorf("doubleFromInt64(%d) reason = %v; want type_mismatch", tc.in, ReasonOf(err))
			}
		}
	}
}

// TestInt32FromFloat64 pins "double → integer exact-fit only".
func TestInt32FromFloat64(t *testing.T) {
	cases := []struct {
		in   float64
		want int32
		ok   bool
	}{
		{5, 5, true},
		{-5, -5, true},
		{5.5, 0, false},
		{math.MinInt32, math.MinInt32, true},
		{math.MaxInt32, math.MaxInt32, true},
		{math.MaxInt32 + 1, 0, false},
		{math.MinInt32 - 1, 0, false},
		{math.NaN(), 0, false},
		{math.Inf(1), 0, false},
		{math.Inf(-1), 0, false},
		{math.Copysign(0, -1), 0, true}, // negative zero is integral
	}
	for _, tc := range cases {
		got, err := int32FromFloat64(tc.in)
		if tc.ok != (err == nil) || (tc.ok && got != tc.want) {
			t.Errorf("int32FromFloat64(%v) = %v, %v; want ok=%v val=%v", tc.in, got, err, tc.ok, tc.want)
		}
	}
}

// TestInt64FromFloat64 covers the JSON longinteger-from-number path.
func TestInt64FromFloat64(t *testing.T) {
	cases := []struct {
		in   float64
		want int64
		ok   bool
	}{
		{1e15, 1000000000000000, true},
		{9007199254740992, 1 << 53, true},
		{1.5, 0, false},
		{maxInt64Float, 0, false},             // 2^63 itself is out of range
		{-maxInt64Float, math.MinInt64, true}, // -2^63 is in range
		{1e300, 0, false},
		{math.NaN(), 0, false},
		{math.Inf(-1), 0, false},
	}
	for _, tc := range cases {
		got, err := int64FromFloat64(tc.in)
		if tc.ok != (err == nil) || (tc.ok && got != tc.want) {
			t.Errorf("int64FromFloat64(%v) = %v, %v; want ok=%v val=%v", tc.in, got, err, tc.ok, tc.want)
		}
	}
}

// TestInt32FromInt64 covers the BSON int64-in-integer-mapping narrow.
func TestInt32FromInt64(t *testing.T) {
	if v, err := int32FromInt64(math.MaxInt32); err != nil || v != math.MaxInt32 {
		t.Errorf("int32FromInt64(MaxInt32) = %v, %v", v, err)
	}
	if _, err := int32FromInt64(math.MaxInt32 + 1); ReasonOf(err) != ReasonTypeMismatch {
		t.Errorf("int32FromInt64(MaxInt32+1) reason = %v; want type_mismatch", ReasonOf(err))
	}
	if _, err := int32FromInt64(math.MinInt32 - 1); err == nil {
		t.Error("int32FromInt64(MinInt32-1) accepted")
	}
}

// TestCheckDouble pins NaN/±Inf rejection.
func TestCheckDouble(t *testing.T) {
	for _, bad := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := checkDouble(bad); ReasonOf(err) != ReasonTypeMismatch {
			t.Errorf("checkDouble(%v) reason = %v; want type_mismatch", bad, ReasonOf(err))
		}
	}
	if v, err := checkDouble(22.5); err != nil || v != 22.5 {
		t.Errorf("checkDouble(22.5) = %v, %v", v, err)
	}
}

// TestDateTimeBounds pins the [MinDateTime, MaxDateTime] window on both the
// epoch-ms and parsed-instant entry points.
func TestDateTimeBounds(t *testing.T) {
	if ts, err := dateTimeFromMillis(minDateTimeMillis); err != nil || !ts.Equal(MinDateTime) {
		t.Errorf("dateTimeFromMillis(min) = %v, %v", ts, err)
	}
	if ts, err := dateTimeFromMillis(maxDateTimeMillis); err != nil || !ts.Equal(MaxDateTime) {
		t.Errorf("dateTimeFromMillis(max) = %v, %v", ts, err)
	}
	for _, bad := range []int64{minDateTimeMillis - 1, maxDateTimeMillis + 1, math.MinInt64, math.MaxInt64} {
		if _, err := dateTimeFromMillis(bad); ReasonOf(err) != ReasonBadTimestamp {
			t.Errorf("dateTimeFromMillis(%d) reason = %v; want bad_timestamp", bad, ReasonOf(err))
		}
	}

	if _, err := checkDateTime(MaxDateTime.Add(time.Millisecond)); ReasonOf(err) != ReasonBadTimestamp {
		t.Errorf("checkDateTime(max+1ms) reason = %v; want bad_timestamp", ReasonOf(err))
	}
	if _, err := checkDateTime(MinDateTime.Add(-time.Nanosecond)); ReasonOf(err) != ReasonBadTimestamp {
		t.Errorf("checkDateTime(min-1ns) reason = %v; want bad_timestamp", ReasonOf(err))
	}
	// Sub-millisecond precision is truncated, not rejected.
	in := time.Date(2026, 6, 10, 12, 34, 56, 789_654_321, time.UTC)
	got, err := checkDateTime(in)
	want := time.Date(2026, 6, 10, 12, 34, 56, 789_000_000, time.UTC)
	if err != nil || !got.Equal(want) {
		t.Errorf("checkDateTime(%v) = %v, %v; want %v", in, got, err, want)
	}
}

// TestCheckString pins the UTF-8 requirement and the 64 KiB string cap.
func TestCheckString(t *testing.T) {
	if _, err := checkString(string([]byte{0xff, 0xfe})); ReasonOf(err) != ReasonTypeMismatch {
		t.Errorf("invalid UTF-8 reason = %v; want type_mismatch", ReasonOf(err))
	}
	atCap := string(make([]byte, MaxStringLen))
	if _, err := checkString(atCap); err != nil {
		t.Errorf("string at cap rejected: %v", err)
	}
	if _, err := checkString(atCap + "x"); ReasonOf(err) != ReasonValueTooLarge {
		t.Errorf("string over cap reason = %v; want value_too_large", ReasonOf(err))
	}
}

// TestCollectArray pins the 1024-element cap and index-annotated errors.
func TestCollectArray(t *testing.T) {
	if _, err := collectArray[int32](MaxArrayLen+1, nil); ReasonOf(err) != ReasonValueTooLarge {
		t.Errorf("oversize array reason = %v; want value_too_large", ReasonOf(err))
	}
	got, err := collectArray[int32](3, func(i int) (Value, error) { return int32(i), nil })
	if err != nil || len(got) != 3 || got[2] != 2 {
		t.Errorf("collectArray = %v, %v", got, err)
	}
	_, err = collectArray[int32](2, func(i int) (Value, error) {
		if i == 1 {
			return nil, rejectf(ReasonTypeMismatch, "boom")
		}
		return int32(0), nil
	})
	var re *RejectError
	if !errors.As(err, &re) || re.Detail != "element 1: boom" {
		t.Errorf("element error not index-annotated: %v", err)
	}
}

// TestRejectReasonLabels freezes the metrics labels and the helper surface.
func TestRejectReasonLabels(t *testing.T) {
	want := map[RejectReason]string{
		ReasonTooLarge:        "too_large",
		ReasonUnknownFormat:   "unknown_format",
		ReasonMalformed:       "malformed",
		ReasonNoValue:         "no_value",
		ReasonBadTimestamp:    "bad_timestamp",
		ReasonTypeMismatch:    "type_mismatch",
		ReasonValueTooLarge:   "value_too_large",
		ReasonBadObject:       "bad_object",
		ReasonUnsetNotAllowed: "unset_not_allowed",
	}
	all := RejectReasons()
	if len(all) != len(want) {
		t.Fatalf("RejectReasons() returned %d reasons; want %d", len(all), len(want))
	}
	for _, r := range all {
		if r.String() != want[r] {
			t.Errorf("reason %d label = %q; want %q", r, r.String(), want[r])
		}
	}
	if ReasonNone.String() != "RejectReason(0)" {
		t.Errorf("ReasonNone label = %q", ReasonNone.String())
	}
	if ReasonOf(nil) != ReasonNone || ReasonOf(errors.New("plain")) != ReasonNone {
		t.Error("ReasonOf must return ReasonNone for nil and untyped errors")
	}
	err := rejectf(ReasonTooLarge, "%d bytes", 7)
	if got := err.Error(); got != "payload rejected (too_large): 7 bytes" {
		t.Errorf("RejectError.Error() = %q", got)
	}
	if got := fmt.Sprintf("%v", FormatBSON); got != "bson" {
		t.Errorf("FormatBSON label = %q", got)
	}
}
