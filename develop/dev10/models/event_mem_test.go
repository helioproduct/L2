package models

import (
	"dev11/structs"
	"testing"
	"time"
)

func TestEventModelCRUD(t *testing.T) {
	A := structs.MakeEventNoId(
		structs.UserID(1),
		time.Date(2019, 9, 9, 0, 0, 0, 0, time.UTC),
	)

	// Create
	model := NewEventModelMemory()
	eA, err := model.Create(A)
	if err != nil {
		t.Fatal("err should be nil", err)
	}

	// SelectByIdx
	esA, err := model.SelectById(eA.GetId())
	if err != nil {
		t.Fatal("err should be nil", err)
	}
	if esA != eA {
		t.Fatalf("got %v; want %v\n", esA, eA)
	}

	// Update
	newDate := time.Date(2023, 9, 9, 0, 0, 0, 0, time.UTC)
	if !eA.SetDate(newDate) {
		t.Fatal("wtf")
	}

	euA, err := model.Update(eA)
	if err != nil {
		t.Fatal("err should be nil", err)
	}
	if euA != eA {
		t.Fatalf("got %v; want %v\n", euA, eA)
	}
	// Delete
	if err := model.Delete(eA.GetId()); err != nil {
		t.Fatal("err should be nil", err)
	}

	if err := model.Delete(eA.GetId()); err == nil {
		t.Fatal("err should be not nil")
	}
}

func eventModelCreateHelper(t *testing.T, m *EventModelMemory, e structs.EventNoId) structs.Event {
	t.Helper()
	eA, err := m.Create(e)
	if err != nil {
		t.Fatal("err should be nil", err)
	}
	return eA
}

func checkEventsSlice(t *testing.T, a, b []structs.Event) {
	t.Helper()
	if len(a) != len(b) {
		t.Fatalf("a != b\na: %v\nb: %v\n", a, b)
	}

	for _, ea := range a {
		var found bool
		for _, eb := range b {
			if ea == eb {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("a != b\na: %v\nb: %v\n", a, b)
		}
	}
}

func TestEventModelSelectBetween(t *testing.T) {
	start := time.Date(2019, 9, 9, 0, 0, 0, 0, time.UTC)
	mid := time.Date(2021, 9, 9, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 9, 9, 0, 0, 0, 0, time.UTC)
	A := structs.MakeEventNoId(
		structs.UserID(1),
		start,
	)
	B := structs.MakeEventNoId(
		structs.UserID(1),
		end,
	)
	C := structs.MakeEventNoId(
		structs.UserID(2),
		mid,
	)

	m := NewEventModelMemory()
	// create
	eB := eventModelCreateHelper(t, m, B)
	eA := eventModelCreateHelper(t, m, A)
	eC := eventModelCreateHelper(t, m, C)

	// find between [mid mid]
	got, err := m.SelectBetweenDates(mid, mid)
	if err != nil {
		t.Fatal("err should be nil")
	}
	checkEventsSlice(t, got, []structs.Event{eC})

	// find between [start mid]
	got, err = m.SelectBetweenDates(start, mid)
	if err != nil {
		t.Fatal("err should be nil")
	}
	checkEventsSlice(t, got, []structs.Event{eA, eC})

	// find between [start end]
	got, err = m.SelectBetweenDates(start, end)
	if err != nil {
		t.Fatal("err should be nil")
	}
	checkEventsSlice(t, got, []structs.Event{eA, eB, eC})
}
