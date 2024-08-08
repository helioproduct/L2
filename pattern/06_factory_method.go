package pattern

/*
	Реализовать паттерн «фабричный метод».
	Объяснить применимость паттерна, его плюсы и минусы,
	а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern
*/

// порождающий
// минусы - большое кол-во параллельных структур
// Привязка к универсальному конструктору

// магазин по продаже серверного оборудования

import (
	"fmt"
)

// Product - интерфейс, определяющий методы, которые должны реализовывать конкретные продукты
type Product2 interface {
	Use() string
}

// ConcreteProductA - конкретная реализация продукта A
type ConcreteProductA struct{}

func (p *ConcreteProductA) Use() string {
	return "Using Product A"
}

// ConcreteProductB - конкретная реализация продукта B
type ConcreteProductB struct{}

func (p *ConcreteProductB) Use() string {
	return "Using Product B"
}

// Creator - интерфейс фабрики, определяющий фабричный метод
type Creator interface {
	CreateProduct(productType string) Product
}

// ConcreteCreator - конкретная реализация фабрики
type ConcreteCreator struct{}

func (c *ConcreteCreator) CreateProduct(productType string) Product {
	if productType == "A" {
		return &ConcreteProductA{}
	} else if productType == "B" {
		return &ConcreteProductB{}
	}
	return nil
}

func main() {
	creator := &ConcreteCreator{}

	productA := creator.CreateProduct("A")
	fmt.Println(productA.Use())

	productB := creator.CreateProduct("B")
	fmt.Println(productB.Use())
}
