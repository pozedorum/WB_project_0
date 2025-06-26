package cache

import (
	"context"
	"log"
	"sync"
	"wb_project_0/internal/db"
	"wb_project_0/internal/models"
)

type Cache struct {
	mu     sync.RWMutex
	orders map[string]models.Order
	db     *db.Database
}

func New(db *db.Database) (*Cache, error) {
	c := &Cache{
		orders: make(map[string]models.Order),
		db:     db,
	}
	err := c.loadFromDB()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Cache) loadFromDB() error {
	orders, err := c.db.GetAllOrders(context.Background())
	if err != nil {
		log.Printf("Failed to load cache from DB: %v", err)
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	for _, order := range orders {
		c.orders[order.OrderUID] = order
	}
	return nil
}

func (c *Cache) Insert(order models.Order) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.db.SaveOrder(context.Background(), order); err != nil {
		log.Printf("Failed to save order to DB: %v", err)
		return err
	}

	c.orders[order.OrderUID] = order
	log.Printf("Order %s cached", order.OrderUID)
	return nil
}

func (c *Cache) Get(uid string) (models.Order, error) {
	c.mu.RLock()
	order, exists := c.orders[uid]
	c.mu.RUnlock()

	if exists {
		return order, nil
	}

	// Если нет в кэше, пробуем загрузить из БД
	orderPtr, err := c.db.GetOrderByUID(context.Background(), uid)
	if err != nil {
		return models.Order{}, err
	}

	c.Insert(*orderPtr) // Добавляем в кэш
	return *orderPtr, nil
}
