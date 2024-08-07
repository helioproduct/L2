package endpoints

import (
	"dev11/logic"
	"dev11/structs"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type EventHTTP struct {
	api *logic.EventAPI
}

func NewEventHTTP(api *logic.EventAPI) *EventHTTP {
	return &EventHTTP{
		api: api,
	}
}

// Формирует и отсылает json документ в w
func (e *EventHTTP) jsonResponse(w http.ResponseWriter, r interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(r)
}

func (e *EventHTTP) jsonResponseBytes(w http.ResponseWriter, b []byte, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(b)
}

// Структура для json сообщения об ошибки
type jsonError struct {
	Error string `json:"error"`
}

// Структура, содержащая один event
type jsonResultEvent struct {
	Result jsonEvent `json:"result"`
}

// Структура, содержащая строковый ответ
type jsonResultString struct {
	Result string `json:"result"`
}

// Структура, содержащая list of events
type jsonResultListOfEvents struct {
	Result []jsonEvent `json:"result"`
}

func (e *EventHTTP) CreateHandle(w http.ResponseWriter, r *http.Request) {
	// Считываем тело запроса
	b, err := io.ReadAll(r.Body)
	if err != nil {
		e.jsonResponse(w, jsonError{err.Error()}, http.StatusBadRequest)
		return
	}
	newe, ok := eventNoIdFromUrlQuery(string(b))
	if !ok {
		e.jsonResponse(w, jsonError{"can't parse"}, http.StatusBadRequest)
		return
	}

	// Бизнес логика
	event, err := e.api.Create(newe)
	if err != nil {
		e.jsonResponse(w, jsonError{err.Error()}, http.StatusInternalServerError)
		return
	}

	// Формируем ответ
	res := jsonResultEvent{
		Result: makeJsonEvent(event),
	}
	e.jsonResponse(w, res, http.StatusOK)
}

func (e *EventHTTP) UpdateHandle(w http.ResponseWriter, r *http.Request) {
	// Считываем тело запроса
	b, err := io.ReadAll(r.Body)
	if err != nil {
		e.jsonResponse(w, jsonError{err.Error()}, http.StatusBadRequest)
		return
	}
	event, ok := eventFromUrlQuery(string(b))
	if !ok {
		e.jsonResponse(w, jsonError{"can't parse"}, http.StatusBadRequest)
		return
	}

	// Бизнес логика
	event, err = e.api.Update(event)
	if err != nil {
		e.jsonResponse(w, jsonError{err.Error()}, http.StatusInternalServerError)
		return
	}

	// Формируем ответ
	res := jsonResultEvent{
		Result: makeJsonEvent(event),
	}
	e.jsonResponse(w, res, http.StatusOK)
}

func (e *EventHTTP) DeleteHandle(w http.ResponseWriter, r *http.Request) {
	// Считываем тело запроса
	b, err := io.ReadAll(r.Body)
	if err != nil {
		e.jsonResponse(w, jsonError{err.Error()}, http.StatusBadRequest)
		return
	}
	event, ok := idFromUrlQuery(structs.EventNoId{}, string(b))
	if !ok {
		e.jsonResponse(w, jsonError{"can't parse"}, http.StatusBadRequest)
		return
	}

	// Бизнес логика
	err = e.api.Delete(event.GetId())
	if err != nil {
		e.jsonResponse(w, jsonError{err.Error()}, http.StatusInternalServerError)
		return
	}

	// Формируем ответ
	res := jsonResultString{
		Result: "deleted",
	}
	e.jsonResponse(w, res, http.StatusOK)
}

type forfunc func(time.Time) ([]structs.Event, error)

func (e *EventHTTP) forFuncHandle(fn forfunc, w http.ResponseWriter, r *http.Request) {
	// Получаем дату
	v := r.URL.Query()
	to_date, ok := freeDateFromUrlValues("to_date", v)
	if !ok {
		e.jsonResponse(w, jsonError{"can't parse"}, http.StatusBadRequest)
		return
	}

	// Бизнес логика
	list, err := fn(to_date)
	if err != nil {
		e.jsonResponse(w, jsonError{err.Error()}, http.StatusInternalServerError)
		return
	}

	// Формируем ответ
	if len(list) == 0 {
		e.jsonResponseBytes(w, []byte(`{"result":[]}`), http.StatusOK)
		return
	}
	res := jsonResultListOfEvents{
		Result: makeSliceJsonEvent(list),
	}
	e.jsonResponse(w, res, http.StatusOK)
}
func (e *EventHTTP) ForDayHandle(w http.ResponseWriter, r *http.Request) {
	e.forFuncHandle(e.api.ForDay, w, r)
}

func (e *EventHTTP) ForWeekHandle(w http.ResponseWriter, r *http.Request) {
	e.forFuncHandle(e.api.ForWeek, w, r)
}

func (e *EventHTTP) ForMonthHandle(w http.ResponseWriter, r *http.Request) {
	e.forFuncHandle(e.api.ForMonth, w, r)
}
