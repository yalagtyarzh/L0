package repository

import "github.com/yalagtyarzh/L0/models"

type DatabaseRepo interface {
	InsertOrder(o models.Order) error
	InsertPaymentByOrderID(pm models.Payment, uid string) error
	InsertDeliveryByOrderID(d models.Delivery, uid string) error
	InsertItem(i models.Item) error
	InsertOrderItems(orderUID, chrtUID string) error
	GetOrders() ([]models.Order, error)
	GetDeliveryByOrderUID(uid string) (models.Delivery, error)
	GetPaymentByOrderUID(uid string) (models.Payment, error)
	GetItemsByOrderUID(uid string) ([]models.Item, error)
}
