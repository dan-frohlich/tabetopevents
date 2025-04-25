package tte

import "testing"

func TestDayTimeCompare(t *testing.T) {
	dt1 := Daytime("Thursday at  1:00 PM")
	dt2 := Daytime("Thursday at  9:00 AM")

	if dt1.Compare(dt2) {
		t.Errorf("[%s] was less than [%s]", dt1, dt2)
	}
	if !dt2.Compare(dt1) {
		t.Errorf("[%s] was greater than [%s]", dt2, dt1)
	}

	dt3 := Daytime("Thursday at  9:00 PM")

	if !dt1.Compare(dt3) {
		t.Errorf("[%s] was greater than [%s]", dt1, dt2)
	}
	if dt3.Compare(dt1) {
		t.Errorf("[%s] was less than [%s]", dt2, dt1)
	}
}
