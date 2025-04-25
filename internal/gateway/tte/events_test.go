package tte

import "testing"

func TestDayTimeCompare(t *testing.T) {
	dt2 := Daytime("Thursday at  9:00 AM")
	dt4 := Daytime("Thursday at 11:00 AM")
	dt1 := Daytime("Thursday at  1:00 PM")
	dt3 := Daytime("Thursday at  9:00 PM")

	if dt1.Compare(dt2) {
		t.Errorf("[%s] was less than [%s]", dt1, dt2)
	}
	if !dt2.Compare(dt1) {
		t.Errorf("[%s] was greater than [%s]", dt2, dt1)
	}

	if !dt1.Compare(dt3) {
		t.Errorf("[%s] was greater than [%s]", dt1, dt3)
	}
	if dt3.Compare(dt1) {
		t.Errorf("[%s] was less than [%s]", dt3, dt1)
	}

	if dt1.Compare(dt4) {
		t.Errorf("[%s] was less than [%s]", dt1, dt4)
	}
	if !dt4.Compare(dt1) {
		t.Errorf("[%s] was greater than [%s]", dt4, dt1)
	}

	if !dt2.Compare(dt4) {
		t.Errorf("[%s] was greater than [%s]", dt2, dt4)
	}
	if dt4.Compare(dt2) {
		t.Errorf("[%s] was less than [%s]", dt4, dt2)
	}
}
