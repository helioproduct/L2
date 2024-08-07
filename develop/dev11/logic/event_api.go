package logic

import (
	"dev11/structs"
	"time"
)

// Интерфейс для работы с моделю event
type IEventsModel interface {
	Create(newe structs.EventNoId) (structs.Event, error)
	SelectById(id structs.EventID) (structs.Event, error)
	SelectBetweenDates(start, end time.Time) ([]structs.Event, error)
	Update(e structs.Event) (structs.Event, error)
	Delete(id structs.EventID) error
}

type EventAPI struct {
	m IEventsModel
}

func NewEventAPI(m IEventsModel) *EventAPI {
	return &EventAPI{
		m: m,
	}
}

func (api *EventAPI) Create(newe structs.EventNoId) (structs.Event, error) {
	return api.m.Create(newe)
}
func (api *EventAPI) Update(e structs.Event) (structs.Event, error) {
	return api.m.Update(e)
}
func (api *EventAPI) Delete(id structs.EventID) error {
	return api.m.Delete(id)
}

// Возвращает список событий, которые имеют дату date
func (api *EventAPI) ForDay(date time.Time) ([]structs.Event, error) {
	return api.m.SelectBetweenDates(date, date)
}

// Возвращает список событий из дипазаона [date-7days, date]
func (api *EventAPI) ForWeek(date time.Time) ([]structs.Event, error) {
	const week = 7 * 24 * time.Hour
	return api.m.SelectBetweenDates(date.Add(-week), date)
}

// Возвращает список событий из дипазаона [date-30days, date]
// Возможно это немного не тот функционал, который требуется, но так проще
func (api *EventAPI) ForMonth(date time.Time) ([]structs.Event, error) {
	const month = 30 * 24 * time.Hour
	return api.m.SelectBetweenDates(date.Add(-month), date)
}
