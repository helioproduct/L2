package endpoints

import (
	"dev11/structs"
	"net/url"
	"strconv"
	"time"
)

func parseIntFromValues(key string, v url.Values) (int, bool) {
	if userId, ok := v[key]; ok {
		if len(userId) != 1 {
			return 0, false
		}
		if userIdint, err := strconv.Atoi(userId[0]); err == nil {
			return userIdint, true
		}
	}
	return 0, false
}

// Парсит user_id в event из url.Values
func userIdFromUrlValues(e structs.EventNoId, v url.Values) (structs.EventNoId, bool) {
	if num, ok := parseIntFromValues("user_id", v); ok {
		if e.SetUserId(structs.UserID(structs.UserID(num))) {
			return e, true
		}
	}

	return e, false
}

func freeDateFromUrlValues(key string, v url.Values) (time.Time, bool) {
	var res time.Time
	if date, ok := v[key]; ok {
		if len(date) != 1 {
			return res, false
		}
		t, err := time.Parse("2006-01-02", date[0])
		if err != nil {
			return res, false
		}
		return t, true
	}
	return res, false
}

// Парсит date в event из url.Values
func dateFromUrlValues(e structs.EventNoId, v url.Values) (structs.EventNoId, bool) {
	if date, ok := v["date"]; ok {
		if len(date) != 1 {
			return e, false
		}
		t, err := time.Parse("2006-01-02", date[0])
		if err != nil {
			return e, false
		}
		if !e.SetDate(t) {
			return e, false
		}
	} else {
		return e, false
	}
	return e, true
}

func eventNoIdFromValues(values url.Values) (structs.EventNoId, bool) {
	var res structs.EventNoId
	// user_id
	res, ok := userIdFromUrlValues(res, values)
	if !ok {
		return res, false
	}
	// date
	res, ok = dateFromUrlValues(res, values)
	if !ok {
		return res, false
	}
	return res, true
}

// Получает eventNoId из urlquery
// Все поля должны быть заполнены:
// user_id=3&date=2019-09-09
func eventNoIdFromUrlQuery(q string) (structs.EventNoId, bool) {
	var res structs.EventNoId
	values, err := url.ParseQuery(q)
	if err != nil {
		return res, false
	}
	return eventNoIdFromValues(values)
}

// Получает id из values
func idFromUrlValues(e structs.EventNoId, v url.Values) (structs.Event, bool) {
	var res structs.Event
	if num, ok := parseIntFromValues("id", v); ok {
		if res, ok := e.MakeEventWithId(structs.EventID(num)); ok {
			return res, true
		}
	}

	return res, false
}

// Получает id из urlquery
func idFromUrlQuery(e structs.EventNoId, q string) (structs.Event, bool) {
	var res structs.Event
	v, err := url.ParseQuery(q)
	if err != nil {
		return res, false
	}
	if num, ok := parseIntFromValues("id", v); ok {
		if res, ok := e.MakeEventWithId(structs.EventID(num)); ok {
			return res, true
		}
	}
	return res, false
}

// Получает event из urlquery
// Все поля должны быть заполнены:
// id=0&user_id=3&date=2019-09-09
func eventFromUrlQuery(q string) (structs.Event, bool) {
	var res structs.Event
	values, err := url.ParseQuery(q)
	if err != nil {
		return res, false
	}

	// EventNoId
	eni, ok := eventNoIdFromValues(values)
	if !ok {
		return res, false
	}

	// id
	res, ok = idFromUrlValues(eni, values)
	if !ok {
		return res, false
	}

	return res, true
}

// json struct for Event
type jsonEvent struct {
	Id structs.EventID `json:"id"`
	jsonEventNoId
}

func makeSliceJsonEvent(l []structs.Event) []jsonEvent {
	var res []jsonEvent
	for _, e := range l {
		res = append(res, makeJsonEvent(e))
	}
	return res
}

func makeJsonEvent(e structs.Event) jsonEvent {
	return jsonEvent{
		Id:            e.GetId(),
		jsonEventNoId: makeJsonEventNoId(e.EventNoId),
	}
}

type jsonEventNoId struct {
	UserID structs.UserID `json:"user_id"`
	Date   string         `json:"date"`
}

func makeJsonEventNoId(e structs.EventNoId) jsonEventNoId {
	return jsonEventNoId{
		UserID: e.GetUserId(),
		Date:   e.GetDate().Format("2006-01-02"),
	}
}
