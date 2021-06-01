package observe

import (
	"testing"
)

func (s *Shirt) UpdateShirtStatusForTest(status string) {
	if s.status != status {
		s.status = status
	}
}

func TestRushToBuy(t *testing.T) {
	c1 := NewCustomers(1)
	c2 := NewCustomers(2)
	c3 := NewCustomers(3)
	s := NewShirt()
	s.Register(c1)
	s.Register(c2)
	s.Register(c3)
	s.UpdateShirtStatusForTest(TimeIsUp)
	s.NotifyObservers()
}
