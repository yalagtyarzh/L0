package repository

import (
	"github.com/yalagtyarzh/L0/internal/models"
)

type DatabaseRepo interface {
	InsertOrder(o models.Order) error
	InsertPaymentByOrderID(pm models.Payment, orderID int) error
	InsertDeliveryByOrderID(d models.Delivery, orderID int) error
	InsertItem(i models.Item) (int, error)
	InsertOrderItems(orderID, chrtID int) error
	GetOrders() ([]models.Order, error)
	GetDeliveryByOrderUID(orderID int) (models.Delivery, error)
	GetPaymentByOrderUID(orderID int) (models.Payment, error)
	GetItemsByOrderUID(orderID int) ([]models.Item, error)
}
