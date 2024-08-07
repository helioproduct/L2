package pattern

import (
	"errors"
	"fmt"
	"time"
)

/*
	Реализовать паттерн «фасад».
	Объяснить применимость паттерна, его плюсы и минусы,
	а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/

// Структурный паттерн
// простой интерфейс к сложной системе
// изоляция клиентов от системы
//

// Плюсы: упрощение сложной системы для пользователей, предоставление
// простого интерфейса для использоавние системы
// Минусы: чрезмерное использование его в разных частях приложения в виде так называемого "супер-класса"

type Product struct {
	Name  string
	Price float64
}

type Shop struct {
	Name     string
	Products []Product
}

// Фасад: подсистемы пользователь, банк, магазин
func (shop Shop) Sell(user User, product string) error {
	println("[Shop] Запрос к пользователю для получение остатка на карте")
	time.Sleep(time.Millisecond * 500)
	err := user.Card.CheckBalance()
	if err != nil {
		return err
	}
	fmt.Printf("[Shop] проверка может ли пользователь (%s) купить товар\n", user.Name)

	for _, prod := range shop.Products {
		if prod.Name != product {
			continue
		}
		if prod.Price > user.GetBalance() {
			return errors.New("[Shop] у пользоватей недостаточно денег для покупки")
		}
	}
	fmt.Printf("[Shop] товар [%s] куплен пользоваталеем", product, user.Name)
	return nil
}

type Card struct {
	Name    string
	Bank    *Bank
	Balance *float64
}

func (card Card) CheckBalance() error {
	println("[Card] запрос в банк для проверки остатка")
	card.Bank.CheckBalance(card.Name)
	time.Sleep(time.Millisecond * 800)
}

type Bank struct {
	Name  string
	Cards []Card
}

func (bank Bank) CheckBalance(cardNumber string) error {
	println("[Bank] получение остатка по карте:", cardNumber)
	time.Sleep(300 * time.Millisecond)

	for _, card := range bank.Cards {
		if card.Name != cardNumber {
			continue
		}
		if *card.Balance <= 0 {
			return errors.New("[Bank] недостаточно средств")
		}
	}
	println("[Bank] Остаток положительный")
}

type User struct {
	Name string
	Card *Card
}

func (user User) GetBalance() float64 {
	return *user.Card.Balance
}

// платежная система
