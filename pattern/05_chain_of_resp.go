package pattern

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

// Плюсы: уменьшает зависимость между клиентами и обработчиками

import (
	"fmt"
)

// Определяем интерфейс обработчика
type Handler interface {
	SetNext(handler Handler)
	Handle(request string)
}

// Базовый обработчик, который хранит следующий обработчик в цепочке
type BaseHandler struct {
	next Handler
}

func (h *BaseHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *BaseHandler) Handle(request string) {
	if h.next != nil {
		h.next.Handle(request)
	}
}

// Конкретный обработчик, который обрабатывает запрос
type ConcreteHandlerA struct {
	BaseHandler
}

func (h *ConcreteHandlerA) Handle(request string) {
	if request == "A" {
		fmt.Println("Handler A обработал запрос")
	} else {
		fmt.Println("Handler A передал запрос дальше")
		h.BaseHandler.Handle(request)
	}
}

// Еще один конкретный обработчик
type ConcreteHandlerB struct {
	BaseHandler
}

func (h *ConcreteHandlerB) Handle(request string) {
	if request == "B" {
		fmt.Println("Handler B обработал запрос")
	} else {
		fmt.Println("Handler B передал запрос дальше")
		h.BaseHandler.Handle(request)
	}
}

func main() {
	handlerA := &ConcreteHandlerA{}
	handlerB := &ConcreteHandlerB{}

	handlerA.SetNext(handlerB)

	// Передаем запросы
	fmt.Println("Передаем запрос 'A':")
	handlerA.Handle("A")

	fmt.Println("\nПередаем запрос 'B':")
	handlerA.Handle("B")

	fmt.Println("\nПередаем запрос 'C':")
	handlerA.Handle("C")
}
