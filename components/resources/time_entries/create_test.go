package time_entries

// Regression tests for review finding #4 on PR #15, asserting the FIXED
// behaviour of hoursToISO8601: correct ISO 8601 output that round-trips through
// github.com/sosodev/duration, including sub-minute precision (no silent loss).

import (
	"math"
	"testing"

	"github.com/sosodev/duration"
)

func TestHoursToISO8601_FormatAndRoundTrip(t *testing.T) {
	cases := []struct {
		hours float64
		want  string
	}{
		{1.0, "PT1H"},
		{1.5, "PT1H30M"},
		{0.25, "PT15M"},
		{0.5, "PT30M"},
		{8.0, "PT8H"},
		{2.75, "PT2H45M"},
		{25.0, "PT25H"}, // flat hours, no day roll-up
	}

	for _, c := range cases {
		got := hoursToISO8601(c.hours)
		if got != c.want {
			t.Errorf("hoursToISO8601(%g) = %q, want %q", c.hours, got, c.want)
		}

		parsed, err := duration.Parse(got)
		if err != nil {
			t.Errorf("library cannot parse output %q: %v", got, err)
			continue
		}
		if rt := parsed.ToTimeDuration().Hours(); math.Abs(rt-c.hours) > 1e-9 {
			t.Errorf("round-trip of %q = %g hours, want %g", got, rt, c.hours)
		}
	}
}

// A positive sub-minute entry must no longer be rounded away to zero.
func TestHoursToISO8601_SubMinutePreserved(t *testing.T) {
	const tiny = 0.001 // 3.6 seconds
	got := hoursToISO8601(tiny)

	parsed, err := duration.Parse(got)
	if err != nil {
		t.Fatalf("library cannot parse %q: %v", got, err)
	}
	if rt := parsed.ToTimeDuration().Hours(); rt <= 0 {
		t.Errorf("sub-minute input %g formatted as %q -> %g hours; expected > 0", tiny, got, rt)
	}
}
