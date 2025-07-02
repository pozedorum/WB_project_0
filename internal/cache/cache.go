package cache

import (
	"context"
	"log"
	"os"
	"sync"
	"wb_project_0/internal/db"
	"wb_project_0/internal/models"
)

type Cache struct {
	mu      sync.RWMutex
	orders  map[string]models.Order
	db      *db.Database
	logger  *log.Logger
	logFile *os.File
}

func New(db *db.Database) (*Cache, error) {

	//logFile, err := os.OpenFile("cache.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	logger := log.New(os.Stdout, "[CACHE] ", log.LstdFlags)
	c := &Cache{
		orders: make(map[string]models.Order),
		db:     db,
		logger: logger,
	}
	err := c.loadFromDB()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Cache) loadFromDB() error {
	c.logger.Printf("getting all orders")
	orders, err := c.db.GetAllOrders(context.Background())
	if err != nil {

		c.logger.Printf("Failed to load cache from DB: %v", err)
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	for _, order := range orders {
		c.orders[order.OrderUID] = order
	}
	c.logger.Printf("getting cache completed")
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
		c.logger.Printf("order with uid:%s recieved from cache", uid)
		return order, nil
	}

	// Если нет в кэше, пробуем загрузить из БД
	orderPtr, err := c.db.GetOrderByUID(context.Background(), uid)
	if err != nil {
		c.logger.Printf("error: order with uid:%s not found", uid)
		return models.Order{}, err
	}
	c.logger.Printf("order with uid:%s recieved from BD", uid)
	err = c.Insert(*orderPtr)
	if err != nil {
		c.logger.Printf("error: cache failed to insert order %s", uid)
	}
	return *orderPtr, nil
}
