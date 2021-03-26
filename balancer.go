package go_payment_balancer

import (
	"errors"
	"fmt"
	"sync"
)

type Balancer struct {
	m        sync.Mutex
	gateways *gateways
}

func NewBalancer() *Balancer {
	return &Balancer{gateways: newGateways()}
}

func (b *Balancer) AddGateway(id interface{}, weight, counter int) {
	defer b.m.Unlock()
	b.m.Lock()
	b.gateways.add(id, weight, counter)
}

func (b *Balancer) RemoveGateway(id interface{}) {
	defer b.m.Unlock()
	b.m.Lock()
	b.gateways.remove(id)
}

func (b *Balancer) IncrementGateway(id interface{}) (err error) {
	defer b.m.Unlock()
	b.m.Lock()
	if gateway := b.gateways.findById(id); gateway != nil {
		gateway.counter++
		return
	}
	err = errors.New("gateway not found")
	return
}

func (b *Balancer) TotalPayments() (total int) {
	defer b.m.Unlock()
	b.m.Lock()
	return b.gateways.totalCount()
}

func (b *Balancer) TotalPaymentsForId(id int) (total int, err error) {
	defer b.m.Unlock()
	b.m.Lock()
	if gateway := b.gateways.findById(id); gateway != nil {
		total = gateway.counter
		return
	}
	err = errors.New("gateway not found")
	return
}

func (b *Balancer) GetGatewayId() (id interface{}, err error) {
	defer b.m.Unlock()
	b.m.Lock()
	if len(*b.gateways) == 0 {
		err = errors.New("no gateways defined")
		return
	}
	id = b.gateways.random().id
	if b.gateways.totalCount() == 0 {
		return
	}
	for _, gateway := range *b.gateways {
		if gateway.counter < (b.gateways.totalCount()/b.gateways.totalWeight())*gateway.weight {
			id = gateway.id
			return
		}
	}
	return
}

func (b *Balancer) Report() (report string) {
	for _, g := range *b.gateways {
		fmt.Println(g)
	}
	report = fmt.Sprintf("Total Payment Count: %d\n", b.gateways.totalCount())
	for _, g := range *b.gateways {
		report += fmt.Sprintf("Count for id %+v = %d\n", g.id, g.counter)
	}
	return
}

// reset gateways
func (b *Balancer) Reset() {
	defer b.m.Unlock()
	b.m.Lock()
	b.gateways = newGateways()
}
