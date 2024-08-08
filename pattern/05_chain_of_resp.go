package pattern

/*
	Реализовать паттерн «цепочка вызовов».
	Объяснить применимость паттерна, его плюсы и минусы,
	а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

// Плюсы:
// 		уменьшает зависимость между клиентами и обработчиками
// 	    Реализует принцип единственной обязанности.
// 		Реализует принцип открытости/закрытости.

//  Минусы:
// 		 Запрос может остаться никем не обработанным.

import (
	"fmt"
)

// Интерфейс для обработчика
type Doctor interface {
	SetNext(doctor Doctor)
	HandleRequest(complaint string)
}

// Базовый обработчик, реализующий общий функционал для всех врачей
type BaseDoctor struct {
	next Doctor
}

func (d *BaseDoctor) SetNext(doctor Doctor) {
	d.next = doctor
}

func (d *BaseDoctor) HandleRequest(complaint string) {
	if d.next != nil {
		d.next.HandleRequest(complaint)
	}
}

// Терапевт, конкретный обработчик
type Therapist struct {
	BaseDoctor
}

func (d *Therapist) HandleRequest(complaint string) {
	if complaint == "простуда" {
		fmt.Println("Терапевт лечит пациента с жалобой на простуду")
	} else {
		fmt.Println("Терапевт передает запрос дальше")
		d.BaseDoctor.HandleRequest(complaint)
	}
}

// Хирург, конкретный обработчик
type Surgeon struct {
	BaseDoctor
}

func (d *Surgeon) HandleRequest(complaint string) {
	if complaint == "аппендицит" {
		fmt.Println("Хирург оперирует пациента с аппендицитом")
	} else {
		fmt.Println("Хирург передает запрос дальше")
		d.BaseDoctor.HandleRequest(complaint)
	}
}

// Специалист, конкретный обработчик
type Specialist struct {
	BaseDoctor
}

func (d *Specialist) HandleRequest(complaint string) {
	if complaint == "заболевание сердца" {
		fmt.Println("Специалист лечит пациента с заболеванием сердца")
	} else {
		fmt.Println("Специалист не может помочь с данной жалобой")
		d.BaseDoctor.HandleRequest(complaint)
	}
}

func main() {
	therapist := &Therapist{}
	surgeon := &Surgeon{}
	specialist := &Specialist{}

	therapist.SetNext(surgeon)
	surgeon.SetNext(specialist)

	// Передаем запросы
	fmt.Println("Пациент жалуется на простуду:")
	therapist.HandleRequest("простуда")

	fmt.Println("\nПациент жалуется на аппендицит:")
	therapist.HandleRequest("аппендицит")

	fmt.Println("\nПациент жалуется на заболевание сердца:")
	therapist.HandleRequest("заболевание сердца")

	fmt.Println("\nПациент жалуется на неизвестную проблему:")
	therapist.HandleRequest("неизвестная проблема")
}
