package pattern

/*
	Реализовать паттерн «стратегия».
	Объяснить применимость паттерна, его плюсы и минусы,
	а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern
*/

import "fmt"

type Item struct {
	Name  string
	price int
}

type PaymentMethod interface {
	Pay(amount float64) string
}

type CreditCard struct {
	name, cardNumber string
}

func (c *CreditCard) Pay(amount float64) string {
	return fmt.Sprintf("Paid %.2f using Credit Card (%s)", amount, c.cardNumber)
}

type PayPal struct {
	email string
}

func (p *PayPal) Pay(amount float64) string {
	return fmt.Sprintf("Paid %.2f using PayPal (%s)", amount, p.email)
}

type Cryptocurrency struct {
	walletAddress string
}

func (c *Cryptocurrency) Pay(amount float64) string {
	return fmt.Sprintf("Paid %.2f using Cryptocurrency (%s)", amount, c.walletAddress)
}

type ShoppingCart struct {
	items         []Item
	paymentMethod PaymentMethod
}

func (s *ShoppingCart) SetPaymentMethod(paymentMethod PaymentMethod) {
	s.paymentMethod = paymentMethod
}

func (s *ShoppingCart) Checkout() string {
	var total float64
	for _, item := range s.items {
		total += item.price
	}
	return s.paymentMethod.Pay(total)
}

func main() {
	shoppingCart := &ShoppingCart{
		items: []Item{
			{"Laptop", 1500},
			{"Smartphone", 1000},
		},
	}

	creditCard := &CreditCard{"Chidozie C. Okafor", "4111-1111-1111-1111"}
	paypal := &PayPal{"chidosiky2015@gmail.com"}
	cryptocurrency := &Cryptocurrency{"0xAbcDe1234FghIjKlMnOp"}

	shoppingCart.SetPaymentMethod(creditCard)
	fmt.Println(shoppingCart.Checkout())

	shoppingCart.SetPaymentMethod(paypal)
	fmt.Println(shoppingCart.Checkout())

	shoppingCart.SetPaymentMethod(cryptocurrency)
	fmt.Println(shoppingCart.Checkout())
}
