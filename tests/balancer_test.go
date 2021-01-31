package tests

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/aliforever/go_payment_balancer"
)

func TestBalancer(t *testing.T) {
	paymentsCount := 5000
	type Gateway struct {
		Id      int
		Title   string
		Weight  int
		Counter int // If there were previous payments you can add initial value with this field
	}
	var gateways = []Gateway{
		{
			Id:      1,
			Title:   "PayPal",
			Weight:  1,
			Counter: 0,
		},
		{
			Id:      2,
			Title:   "MasterCard",
			Weight:  2,
			Counter: 0,
		},
	}
	// force tag
	shouldIncrement := []bool{true, false} // This is to mimic user behavior, some might cancel payment
	b := go_payment_balancer.NewBalancer()
	for _, gateway := range gateways {
		b.AddGateway(gateway.Id, gateway.Weight, gateway.Counter)
	}
	i := 1
	for true {
		g, err := b.GetGatewayId()
		if err != nil {
			fmt.Println(err)
			return
		}
		if shouldIncrement[rand.Intn(len(shouldIncrement))] {
			b.IncrementGateway(g)
			i++
		}
		if b.TotalPayments() >= paymentsCount {
			break
		}
	}
	fmt.Println(b.Report())
}
