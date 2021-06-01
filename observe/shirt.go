package observe

import (
	"fmt"
	"sync"
)

const (
	TimeIsUp  = "time is up"
	IsEnd     = "is end"
	NotAtTime = "not at time"
)

//被观察者
type Shirt struct {
	sync.Mutex
	customers []Observer
	status    string
	name      string
}

func NewShirt() *Shirt {
	return &Shirt{status: NotAtTime, name: "shirt"}
}

func (s *Shirt) Register(o Observer) {
	s.Lock()
	defer s.Unlock()
	s.customers = append(s.customers, o)
	fmt.Printf("[%s] registered a new customer with ID[%d]\n", s.name, o.GetId())
}

func (s *Shirt) Deregister(o Observer) error {
	s.Lock()
	defer s.Unlock()
	var index int
	var found bool
	id := o.GetId()
	for i, c := range s.customers {
		if c.GetId() == id {
			index = i
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("Customer %d not found\n", id)
	}
	s.customers = append(s.customers[:index], s.customers[index+1:]...)
	fmt.Printf("Removed the customer with ID[%d]\n", id)
	return nil
}

func (s *Shirt) NotifyObservers() {
	s.Lock()
	defer s.Unlock()
	wg := sync.WaitGroup{}
	for _, c := range s.customers {
		wg.Add(1)
		go func(c Observer) {
			defer wg.Done()
			c.Update(s.name, s.status)
		}(c)
	}
	wg.Wait()
	fmt.Println("Finished notify customers")
}
