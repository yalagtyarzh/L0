package repocache

import (
	"sync"

	"github.com/yalagtyarzh/L0/internal/models"
	"github.com/yalagtyarzh/L0/internal/repository"
)

// Cache for holding data about orders
type Cache struct {
	orders map[string]models.Order
	mutex  sync.RWMutex
}

// NewCache returns empty orders cache object pointer
func NewCache() *Cache {
	return &Cache{
		orders: map[string]models.Order{},
		mutex:  sync.RWMutex{},
	}
}

// Recover recovers orders cache from database
func (c *Cache) Recover(repo repository.DatabaseRepo) error {
	o, err := repo.GetOrders()
	if err != nil {
		return err
	}

	orders := map[string]models.Order{}
	for _, order := range o {
		orders[order.OrderUID] = order
	}

	c.orders = orders

	return nil
}

// Load returns order by key
func (c *Cache) Load(key string) (models.Order, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	val, ok := c.orders[key]
	return val, ok
}

// Store stores new order with key into cache
func (c *Cache) Store(key string, value models.Order) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.orders[key] = value
}
