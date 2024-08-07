package endpoints

import (
	"bytes"
	"dev11/logic"
	"dev11/models"
	"dev11/structs"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func buildEventHTTP() *EventHTTP {
	return NewEventHTTP(logic.NewEventAPI(models.NewEventModelMemory()))
}

func checkStatusBody(t *testing.T, method, url, body string, handler http.HandlerFunc, wantCode int, wantBody string) {
	t.Helper()

	// На самом деле пока не важно какого типа запрос
	req, err := http.NewRequest(method, url, bytes.NewReader([]byte(body)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Проверяем status code
	if status := rr.Code; status != wantCode {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Провермя тело
	if rr.Body.String() != wantBody {
		t.Errorf("handler returned unexpected body: \ngot  %v\nwant %v",
			rr.Body, wantBody)
	}
}
func TestCreate(t *testing.T) {
	e := buildEventHTTP()
	body := `user_id=3&date=2019-09-09`

	wantStatusCode := http.StatusOK
	wantBody := `{"result":{"id":0,"user_id":3,"date":"2019-09-09"}}`

	checkStatusBody(t, "", "", body, e.CreateHandle, wantStatusCode, wantBody+"\n")
}

func TestUpdate(t *testing.T) {
	e := buildEventHTTP()
	// Добавим событие
	_, _ = e.api.Create(structs.MakeEventNoId(
		1,
		time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
	))
	// Изменим пользователя и дату
	body := `id=0&user_id=2&date=2020-01-01`

	wantStatusCode := http.StatusOK
	wantBody := `{"result":{"id":0,"user_id":2,"date":"2020-01-01"}}`

	checkStatusBody(t, "", "", body, e.UpdateHandle, wantStatusCode, wantBody+"\n")
}

func TestDelete(t *testing.T) {
	e := buildEventHTTP()
	// Добавим событие
	_, _ = e.api.Create(structs.MakeEventNoId(
		1,
		time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
	))
	// Удалим событие
	body := `id=0`

	wantStatusCode := http.StatusOK
	wantBody := `{"result":"deleted"}`

	checkStatusBody(t, "", "", body, e.DeleteHandle, wantStatusCode, wantBody+"\n")
}

func TestForDay(t *testing.T) {
	e := buildEventHTTP()
	// Добавим события
	_, _ = e.api.Create(structs.MakeEventNoId(
		1,
		time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
	))
	// Получим данные
	url := `?to_date=2019-01-01`

	wantStatusCode := http.StatusOK
	wantBody := `{"result":[{"id":0,"user_id":1,"date":"2019-01-01"}]}`

	checkStatusBody(t, "GET", url, "", e.ForDayHandle, wantStatusCode, wantBody+"\n")
}

func TestForWeek(t *testing.T) {
	e := buildEventHTTP()
	// Добавим события
	_, _ = e.api.Create(structs.MakeEventNoId(
		1,
		time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
	))
	_, _ = e.api.Create(structs.MakeEventNoId(
		1,
		time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC),
	))
	// Получим данные
	url := `?to_date=2019-01-07`

	wantStatusCode := http.StatusOK
	wantBody := `{"result":[{"id":0,"user_id":1,"date":"2019-01-01"},{"id":1,"user_id":1,"date":"2019-01-07"}]}`

	checkStatusBody(t, "GET", url, "", e.ForWeekHandle, wantStatusCode, wantBody+"\n")
}

func TestForMonth(t *testing.T) {
	e := buildEventHTTP()
	// Добавим события
	_, _ = e.api.Create(structs.MakeEventNoId(
		1,
		time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
	))
	_, _ = e.api.Create(structs.MakeEventNoId(
		1,
		time.Date(2019, 1, 30, 0, 0, 0, 0, time.UTC),
	))
	// Получим данные
	url := `?to_date=2019-01-30`

	wantStatusCode := http.StatusOK
	wantBody := `{"result":[{"id":0,"user_id":1,"date":"2019-01-01"},{"id":1,"user_id":1,"date":"2019-01-30"}]}`

	checkStatusBody(t, "GET", url, "", e.ForMonthHandle, wantStatusCode, wantBody+"\n")
}
