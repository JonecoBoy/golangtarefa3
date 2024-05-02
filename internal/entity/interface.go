package entity

type OrderRepositoryInterface interface {
	Save(order *Order) error
	GetOrder(id string) (*Order, error)
	List() ([]*Order, error)
	// GetTotal() (int, error)
}
