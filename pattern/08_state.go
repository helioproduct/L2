package pattern

/*
	Реализовать паттерн «состояние».
	Объяснить применимость паттерна,
	его плюсы и минусы, а также реальные
	примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern
*/

import (
	"fmt"
)

// State - интерфейс состояния
type State interface {
	Handle(context *Context)
	String() string
}

// Context - контекст, который хранит текущее состояние
type Context struct {
	state State
}

// NewContext создает новый контекст с заданным начальным состоянием
func NewContext(initialState State) *Context {
	return &Context{state: initialState}
}

// SetState устанавливает новое состояние контекста
func (c *Context) SetState(state State) {
	c.state = state
}

// Request вызывает обработчик текущего состояния
func (c *Context) Request() {
	c.state.Handle(c)
}

// ConcreteStateA - конкретное состояние A
type ConcreteStateA struct{}

func (s *ConcreteStateA) Handle(context *Context) {
	fmt.Println("State A handling request.")
	context.SetState(&ConcreteStateB{})
}

func (s *ConcreteStateA) String() string {
	return "State A"
}

// ConcreteStateB - конкретное состояние B
type ConcreteStateB struct{}

func (s *ConcreteStateB) Handle(context *Context) {
	fmt.Println("State B handling request.")
	context.SetState(&ConcreteStateA{})
}

func (s *ConcreteStateB) String() string {
	return "State B"
}

func main() {
	// Создаем контекст с начальными состоянием ConcreteStateA
	context := NewContext(&ConcreteStateA{})

	// Выполняем запросы и меняем состояния
	for i := 0; i < 5; i++ {
		fmt.Printf("Current state: %s\n", context.state)
		context.Request()
	}
}
