package structs

import (
	"testing"
	"time"
)

func TestEventGetters(t *testing.T) {
	et := time.Date(2001, 05, 13, 0, 0, 0, 0, time.UTC)
	e := Event{
		id: 1,
		EventNoId: EventNoId{
			userId: 2,
			date:   et,
		},
	}

	if e.GetId() != e.id {
		t.Fatalf("got: %v want: %v", e.GetId(), e.id)
	}
	if e.GetUserId() != e.userId {
		t.Fatalf("got: %v want: %v", e.GetUserId(), e.userId)
	}
	if e.GetDate() != e.date {
		t.Fatalf("got: %v want: %v", e.GetDate(), e.date)
	}
}

func TestEventSetters(t *testing.T) {
	et := time.Date(2001, 05, 13, 0, 0, 0, 0, time.UTC)
	e := Event{
		id: 1,
		EventNoId: EventNoId{
			userId: 2,
			date:   et,
		},
	}

	// Date Мы берем только год, день, месяц
	newt := time.Date(2020, 05, 13, 23, 56, 23, 0, time.UTC)
	want := time.Date(2020, 05, 13, 0, 0, 0, 0, time.UTC)
	if !e.SetDate(newt) {
		t.Fatal("wtf")
	}
	if e.GetDate() != want {
		t.Fatalf("got: %v want: %v", e.GetDate(), want)
	}

	// User
	wasu := e.GetUserId()
	if e.SetUserId(-5) {
		t.Fatal("should be pos number (not ok)")
	}
	if e.GetUserId() != wasu {
		t.Fatal("user changed")
	}
	if !e.SetUserId(5) {
		t.Fatal("should be ok")
	}
}
