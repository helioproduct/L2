package pattern

import "fmt"

// Product представляет сложный объект
type Product struct {
	part1 string
	part2 string
	part3 string
}

// String возвращает строковое представление продукта
func (p *Product) String() string {
	return fmt.Sprintf("Product [part1=%s, part2=%s, part3=%s]", p.part1, p.part2, p.part3)
}

// Builder интерфейс определяет шаги для построения продукта
type Builder interface {
	SetPart1(part1 string) Builder
	SetPart2(part2 string) Builder
	SetPart3(part3 string) Builder
	Build() *Product
}

// ConcreteBuilder конкретная реализация Builder
type ConcreteBuilder struct {
	product *Product
}

// NewConcreteBuilder создаёт новый ConcreteBuilder
func NewConcreteBuilder() *ConcreteBuilder {
	return &ConcreteBuilder{product: &Product{}}
}

// SetPart1 задаёт часть 1 продукта
func (b *ConcreteBuilder) SetPart1(part1 string) Builder {
	b.product.part1 = part1
	return b
}

// SetPart2 задаёт часть 2 продукта
func (b *ConcreteBuilder) SetPart2(part2 string) Builder {
	b.product.part2 = part2
	return b
}

// SetPart3 задаёт часть 3 продукта
func (b *ConcreteBuilder) SetPart3(part3 string) Builder {
	b.product.part3 = part3
	return b
}

// Build возвращает построенный продукт
func (b *ConcreteBuilder) Build() *Product {
	return b.product
}

// Director направляет процесс построения
type Director struct {
	builder Builder
}

// NewDirector создаёт нового Director
func NewDirector(builder Builder) *Director {
	return &Director{builder: builder}
}

// Construct определяет порядок создания продукта
func (d *Director) Construct() *Product {
	return d.builder.SetPart1("Part 1").SetPart2("Part 2").SetPart3("Part 3").Build()
}

func main() {
	builder := NewConcreteBuilder()
	director := NewDirector(builder)
	product := director.Construct()
	fmt.Println(product)
}
