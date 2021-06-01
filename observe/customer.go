package observe

import "fmt"

//观察者
type customer struct {
	ID             int
	wantItemStatus string
}

func NewCustomers(id int) *customer {
	return &customer{ID: id}
}

func (c *customer) Update(name, status string) {
	c.wantItemStatus = status
	fmt.Printf("Update: hi customer %d, the item[%s] you want is [%v] now\n", c.ID, name, c.wantItemStatus)
}

func (c *customer) GetId() int {
	return c.ID
}
