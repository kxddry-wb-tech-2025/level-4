package cache

import (
	"container/list"
	c "context"
	"l0/internal/models"
	"l0/internal/storage"
	"sync"
	"time"
)

type cacheEntry struct {
	order   models.Order
	time    time.Time
	lruElem *list.Element
}

// Cache uses TTL + LRU for cache invalidation.
type Cache struct {
	mp       map[string]*cacheEntry
	mu       *sync.Mutex
	ttl      time.Duration
	stopChan chan struct{}
	limit    int
	lru      *list.List
}

// NewCache creates cache
func NewCache(ttl time.Duration, limit int) *Cache {
	cc := &Cache{
		mp:       make(map[string]*cacheEntry),
		mu:       new(sync.Mutex),
		ttl:      ttl,
		stopChan: make(chan struct{}),
		limit:    limit,
		lru:      list.New(),
	}

	go cc.removeExpired()
	return cc
}

// SaveOrder saves
func (c *Cache) SaveOrder(ctx c.Context, order *models.Order) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, ok := c.mp[order.OrderUID]; ok {
		entry.order = *order
		entry.time = time.Now()
		c.lru.MoveToFront(entry.lruElem)
		return nil
	}

	elem := c.lru.PushFront(order.OrderUID)
	c.mp[order.OrderUID] = &cacheEntry{
		order:   *order,
		time:    time.Now(),
		lruElem: elem,
	}

	if c.lru.Len() > c.limit {
		c.removeLRU()
	}

	return nil
}

func (c *Cache) removeLRU() {
	back := c.lru.Back()
	if back == nil {
		return
	}
	orderID := back.Value.(string)
	c.remove(orderID)
}

func (c *Cache) remove(orderID string) {
	entry, ok := c.mp[orderID]
	if !ok {
		return
	}
	c.lru.Remove(entry.lruElem)
	delete(c.mp, orderID)
}

// GetOrder gets an order from the cache
func (c *Cache) GetOrder(ctx c.Context, orderID string) (*models.Order, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.mp[orderID]
	if !ok {
		return nil, storage.ErrOrderNotFound
	}

	// cache invalidation
	if time.Since(entry.time) > c.ttl {
		c.remove(orderID)
		return nil, storage.ErrOrderNotFound
	}

	c.lru.MoveToFront(entry.lruElem)

	return &entry.order, nil
}

// LoadOrders loads orders provided
func (c *Cache) LoadOrders(ctx c.Context, orders []*models.Order) error {
	for _, order := range orders {
		err := c.SaveOrder(ctx, order)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cache) removeExpired() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			c.mu.Lock()
			for id, entry := range c.mp {
				if now.Sub(entry.time) > c.ttl {
					c.remove(id)
				}
			}
			c.mu.Unlock()
		case <-c.stopChan:
			return
		}

	}
}

// Stop stops.
func (c *Cache) Stop() {
	close(c.stopChan)
}
