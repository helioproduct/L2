package models

import (
	"dev11/structs"
	"errors"
	"sync"
	"time"
)

type EventModelMemory struct {
	lock   sync.RWMutex
	events []structs.Event
	freeId structs.EventID
}

func NewEventModelMemory() *EventModelMemory {
	return &EventModelMemory{}
}

func (m *EventModelMemory) findidx(id structs.EventID) (int, bool) {
	for i, v := range m.events {
		if v.GetId() == id {
			return i, true
		}
	}
	return 0, false
}

func (m *EventModelMemory) Create(newe structs.EventNoId) (structs.Event, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	// В качестве нового id берем просто следующий элемент
	newEv, _ := newe.MakeEventWithId(m.freeId)
	m.freeId++
	m.events = append(m.events, newEv)
	return newEv, nil
}

func (m *EventModelMemory) SelectById(id structs.EventID) (structs.Event, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if idx, ok := m.findidx(id); ok {
		return m.events[idx], nil
	}
	return structs.Event{}, errors.New("no such element id")
}

func (m *EventModelMemory) SelectBetweenDates(start, end time.Time) ([]structs.Event, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if end.Before(start) {
		return nil, errors.New("end before start")
	}

	var res []structs.Event
	for _, el := range m.events {
		edata := el.GetDate()
		// [start, end] -> !before(start) && !after(end)
		if !edata.Before(start) && !edata.After(end) {
			res = append(res, el)
		}
	}
	return res, nil
}

func (m *EventModelMemory) Update(e structs.Event) (structs.Event, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if idx, ok := m.findidx(e.GetId()); ok {
		m.events[idx] = e
		return e, nil
	}
	return structs.Event{}, errors.New("no such element id")
}

func (m *EventModelMemory) Delete(id structs.EventID) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if idx, ok := m.findidx(id); ok {
		back := len(m.events) - 1
		m.events[idx], m.events[back] = m.events[back], m.events[idx]
		m.events = m.events[:back]
		return nil
	}
	return errors.New("no such element id")
}
