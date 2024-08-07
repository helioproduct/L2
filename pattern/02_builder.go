package pattern

import "fmt"

/*
	Реализовать паттерн «строитель».
	Объяснить применимость паттерна, его плюсы и минусы,
	а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/

// Порождающий паттерн
// сложный объект = шаги

// Плюсы: создание объекта пошагово
// Минусы: Доп классы

type Collector interface {
	SetCore()
	SetBrand()
	SetMem()
	SetGPU()
	GetComputer()
}

type Computer struct {
	Core   int
	Brand  string
	Memory int
	GPU    int
}

func (pc *Computer) Print() {
	fmt.Printf("Core: [%d] Mem: [%d] GPU: [%d]\n", pc.Core, pc.Memory, pc.GPU)
}

const (
	AsusCollectorType = "asus"
	HpCollectorType   = "hp"
)

type AsusCollector struct {
	Core   int
	Brand  string
	Memory int
	GPU    int
}


func GetCollector(collectorType string) Collector {
	switch collectorType {
	default:
		return nil
	case AsusCollectorType:
		
	case: HpcollectorType:

	}
}
