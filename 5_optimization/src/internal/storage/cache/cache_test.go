package cache_test

import (
	"context"
	"l0/internal/storage/cache"
	"testing"
	"time"

	"l0/internal/models"
	"l0/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestOrder(id string) *models.Order {
	return &models.Order{
		OrderUID:          id,
		TrackNumber:       "im alive!!!",
		Entry:             "",
		Delivery:          models.Delivery{},
		Payment:           models.Payment{},
		Items:             []models.Item{},
		Locale:            "",
		InternalSignature: "",
		CustomerID:        "",
		DeliveryService:   "",
		ShardKey:          "",
		SmID:              0,
		DateCreated:       "",
		OofShard:          "",
	}
}

func TestCache_SaveAndGetOrder(t *testing.T) {
	c := cache.NewCache(5*time.Minute, 10)
	defer c.Stop()

	ctx := context.Background()
	order := newTestOrder("123")

	err := c.SaveOrder(ctx, order)
	require.NoError(t, err)

	got, err := c.GetOrder(ctx, "123")
	require.NoError(t, err)
	assert.Equal(t, order.OrderUID, got.OrderUID)
}

func TestCache_GetOrder_NotFound(t *testing.T) {
	c := cache.NewCache(5*time.Minute, 10)
	defer c.Stop()

	ctx := context.Background()
	_, err := c.GetOrder(ctx, "nonexistent")
	assert.ErrorIs(t, err, storage.ErrOrderNotFound)
}

func TestCache_TTLExpiration(t *testing.T) {
	c := cache.NewCache(10*time.Millisecond, 10)
	defer c.Stop()

	ctx := context.Background()
	order := newTestOrder("expire")

	err := c.SaveOrder(ctx, order)
	require.NoError(t, err)

	time.Sleep(20 * time.Millisecond)

	_, err = c.GetOrder(ctx, "expire")
	assert.ErrorIs(t, err, storage.ErrOrderNotFound)
}

func TestCache_LRUEviction(t *testing.T) {
	c := cache.NewCache(5*time.Minute, 3) // small limit for testing
	defer c.Stop()

	ctx := context.Background()

	// Save orders: 1,2,3
	for i := 1; i <= 3; i++ {
		err := c.SaveOrder(ctx, newTestOrder(string(rune('0'+i))))
		require.NoError(t, err)
	}

	// Access order 1 to make it most recently used
	_, err := c.GetOrder(ctx, "1")
	require.NoError(t, err)

	// Save order 4, should evict LRU (which is 2 now)
	err = c.SaveOrder(ctx, newTestOrder("4"))
	require.NoError(t, err)

	// 2 should be gone
	_, err = c.GetOrder(ctx, "2")
	assert.ErrorIs(t, err, storage.ErrOrderNotFound)

	// 1, 3, and 4 should be present
	for _, id := range []string{"1", "3", "4"} {
		_, err := c.GetOrder(ctx, id)
		assert.NoError(t, err)
	}
}

func TestCache_LoadOrders(t *testing.T) {
	c := cache.NewCache(5*time.Minute, 10)
	defer c.Stop()

	ctx := context.Background()

	orders := []*models.Order{
		newTestOrder("a"),
		newTestOrder("b"),
		newTestOrder("c"),
	}

	err := c.LoadOrders(ctx, orders)
	require.NoError(t, err)

	for _, o := range orders {
		got, err := c.GetOrder(ctx, o.OrderUID)
		require.NoError(t, err)
		assert.Equal(t, o.OrderUID, got.OrderUID)
	}
}
