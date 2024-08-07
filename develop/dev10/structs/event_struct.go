package structs

import "time"

type UserID int // TODO убрать в user.go

type EventID int
type Event struct {
	id EventID
	EventNoId
}

func (e *Event) GetId() EventID {
	return e.id
}

// Event без поля ID, используется для создания записей
type EventNoId struct {
	userId UserID
	date   time.Time
}

// Конструктор для EventNoId
func MakeEventNoId(userId UserID, date time.Time) EventNoId {
	return EventNoId{
		userId: userId,
		date:   date,
	}
}

func (e *EventNoId) GetUserId() UserID {
	return e.userId
}

func (e *EventNoId) GetDate() time.Time {
	return e.date
}

// Пытаемся установить новый userid, если он не корректный
// то возвращаем false
func (e *EventNoId) SetUserId(id UserID) bool {
	if id < 0 {
		return false
	}
	e.userId = id
	return true
}

// Пытаемся установить дату нвоого события.
// При этом задается только дата. Время обнуляется
func (e *EventNoId) SetDate(date time.Time) bool {
	utcDate := date.UTC()
	newDate := time.Date(utcDate.Year(), utcDate.Month(), utcDate.Day(), 0, 0, 0, 0, time.UTC)
	e.date = newDate
	return true
}

// Задаем id событию без id, возвращаем событие с id
func (e EventNoId) MakeEventWithId(id EventID) (Event, bool) {
	if id < 0 {
		return Event{}, false
	}

	return Event{
		id:        id,
		EventNoId: e,
	}, true
}
