package go_payment_balancer

import (
	"errors"
	"math/rand"
	"sort"
)

type gateway struct {
	id      interface{}
	weight  int
	counter int
}

type gateways []*gateway

func newGateways() *gateways {
	return &gateways{}
}

func (g *gateways) add(id interface{}, weight int) error {
	for _, gateway := range *g {
		if gateway.id == id {
			return errors.New("error: id exists")
		}
	}
	*g = append(*g, &gateway{
		id:     id,
		weight: weight,
	})
	g.sortByWeight()
	return nil
}

func (g *gateways) remove(id interface{}) {
	ng := &gateways{}
	for _, gId := range *g {
		if gId.id == id {
			continue
		}
		ng.add(gId.id, gId.weight)
	}
	*g = *ng
}

func (g *gateways) totalWeight() (total int) {
	for _, gateway := range *g {
		total += gateway.weight
	}
	return
}

func (g *gateways) totalCount() (total int) {
	for _, gateway := range *g {
		total += gateway.counter
	}
	return
}

func (g *gateways) findById(id interface{}) *gateway {
	for _, gateway := range *g {
		if gateway.id == id {
			return gateway
		}
	}
	return nil
}

func (g *gateways) sortByWeight() {
	sort.Slice(*g, func(i, j int) bool {
		return (*g)[i].weight > (*g)[j].weight
	})
}

func (g *gateways) random() *gateway {
	return (*g)[rand.Intn(len(*g))]
}
