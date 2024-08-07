package logic

import (
	"dev11/models"
	"dev11/structs"
	"testing"
	"time"
)

func eventAPIMemoryModel() *EventAPI {
	m := models.NewEventModelMemory()
	return NewEventAPI(m)
}

func TestEventAPICreate(t *testing.T) {
	api := eventAPIMemoryModel()

	e := structs.MakeEventNoId(
		1,
		time.Date(2019, 10, 10, 0, 0, 0, 0, time.UTC),
	)
	_, err := api.Create(e)
	if err != nil {
		t.Fatal("err should be nil")
	}
}

func TestEventAPIUpdateDelete(t *testing.T) {
	api := eventAPIMemoryModel()

	newe := structs.MakeEventNoId(
		1,
		time.Date(2019, 10, 10, 0, 0, 0, 0, time.UTC),
	)
	ea, err := api.Create(newe)
	if err != nil {
		t.Fatal("err should be nil")
	}

	// update
	newu := structs.UserID(10)
	if !ea.SetUserId(newu) {
		t.Fatal("wtf")
	}
	ea, err = api.Update(ea)
	if err != nil {
		t.Fatal("err should be nil")
	}
	if ea.GetUserId() != newu {
		t.Fatalf("got: %v want: %v", ea.GetUserId(), newu)
	}

	// delete
	if err := api.Delete(ea.GetId()); err != nil {
		t.Fatal("err should be nil")
	}
	if err := api.m.Delete(ea.GetId()); err == nil {
		t.Fatal("err should be not nil")
	}
}

func TestEventAPIForDay(t *testing.T) {
	api := eventAPIMemoryModel()

	newe := structs.MakeEventNoId(
		1,
		time.Date(2019, 10, 10, 0, 0, 0, 0, time.UTC),
	)
	ea, err := api.Create(newe)
	if err != nil {
		t.Fatal("err should be nil")
	}

	// forday
	toDate := ea.GetDate()
	events, err := api.ForDay(toDate)
	if err != nil {
		t.Fatalf("err should be nil")
	}
	if len(events) != 1 || events[0] != ea {
		t.Fatalf("got %v; want [%v]", events, ea)
	}

	// forday empty
	events, err = api.ForDay(time.Date(2010, 0, 0, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("err should be nil")
	}
	if len(events) != 0 {
		t.Fatal("should be [] or nil")
	}
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

func TestEventAPIForWeek(t *testing.T) {
	api := eventAPIMemoryModel()

	beforeStart := time.Date(2019, 10, 9, 0, 0, 0, 0, time.UTC)
	start := time.Date(2019, 10, 10, 0, 0, 0, 0, time.UTC)
	mid := time.Date(2019, 10, 14, 0, 0, 0, 0, time.UTC)
	end := time.Date(2019, 10, 17, 0, 0, 0, 0, time.UTC)

	e0, _ := api.Create(structs.MakeEventNoId(1, beforeStart))
	ea, _ := api.Create(structs.MakeEventNoId(1, start))
	eb, _ := api.Create(structs.MakeEventNoId(1, mid))
	ec, _ := api.Create(structs.MakeEventNoId(1, end))

	// forweek
	events, err := api.ForWeek(beforeStart)
	if err != nil {
		t.Fatalf("err should be nil")
	}
	checkEventsSlice(t, events, []structs.Event{e0})

	events, err = api.ForWeek(start)
	if err != nil {
		t.Fatalf("err should be nil")
	}
	checkEventsSlice(t, events, []structs.Event{e0, ea})

	events, err = api.ForWeek(mid)
	if err != nil {
		t.Fatalf("err should be nil")
	}
	checkEventsSlice(t, events, []structs.Event{e0, ea, eb})

	events, err = api.ForWeek(end)
	if err != nil {
		t.Fatalf("err should be nil")
	}
	checkEventsSlice(t, events, []structs.Event{ea, eb, ec})
}

func TestEventAPIForMonth(t *testing.T) {
	api := eventAPIMemoryModel()

	beforeStart := time.Date(2018, 12, 31, 0, 0, 0, 0, time.UTC)
	start := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	mid := time.Date(2019, 1, 20, 0, 0, 0, 0, time.UTC)
	end := time.Date(2019, 1, 31, 0, 0, 0, 0, time.UTC)

	e0, _ := api.Create(structs.MakeEventNoId(1, beforeStart))
	ea, _ := api.Create(structs.MakeEventNoId(1, start))
	eb, _ := api.Create(structs.MakeEventNoId(1, mid))
	ec, _ := api.Create(structs.MakeEventNoId(1, end))

	// formonth
	events, err := api.ForMonth(beforeStart)
	if err != nil {
		t.Fatalf("err should be nil")
	}
	checkEventsSlice(t, events, []structs.Event{e0})

	events, err = api.ForMonth(start)
	if err != nil {
		t.Fatalf("err should be nil")
	}
	checkEventsSlice(t, events, []structs.Event{e0, ea})

	events, err = api.ForMonth(mid)
	if err != nil {
		t.Fatalf("err should be nil")
	}
	checkEventsSlice(t, events, []structs.Event{e0, ea, eb})

	events, err = api.ForMonth(end)
	if err != nil {
		t.Fatalf("err should be nil")
	}
	checkEventsSlice(t, events, []structs.Event{ea, eb, ec})
}
