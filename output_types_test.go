package main

import "testing"

func TestConvertTimeUsesOffset(t *testing.T) {
    const timestamp uint = REFERENCE_TIME + 5*SECONDS_IN_HOUR
    const offset int = 2 * int(SECONDS_IN_HOUR)

    _, hour := convertTime(timestamp, 0)
    if hour != 5 {
        t.Errorf("Expected hour 5 with no offset, got %v", hour)
    }

    _, hour = convertTime(timestamp, offset)
    if hour != 7 {
        t.Errorf("Expected hour 7 with offset, got %v", hour)
    }
}
