package cache

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"l0/internal/models"
)

func makeOrder(id string) *models.Order {
	return &models.Order{
		OrderUID:    id,
		TrackNumber: "TRACK" + id,
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "John Doe",
			Phone:   "+12345678901",
			Zip:     "123456",
			City:    "City",
			Address: "Addr",
			Region:  "RG",
			Email:   "john@example.com",
		},
		Payment: models.Payment{
			Transaction:  "txn-" + id,
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       100,
			PaymentDT:    time.Now().Unix(),
			Bank:         "Sber",
			DeliveryCost: 10,
			GoodsTotal:   100,
		},
		Items: []models.Item{{
			ChrtID:      1,
			TrackNumber: "TRACK" + id,
			Price:       100,
			RID:         "rid-" + id,
			Name:        "Item",
			Size:        "M",
			TotalPrice:  100,
			NmID:        1,
			Brand:       "WB",
			Status:      10,
		}},
		Locale:          "en",
		CustomerID:      "cust",
		DeliveryService: "meest",
		ShardKey:        "1",
		SmID:            1,
		DateCreated:     "2021-11-26T06:22:19Z",
		OofShard:        "1",
	}
}

func BenchmarkCache_SaveOrder(b *testing.B) {
	c := NewCache(5*time.Minute, 1_000_000)
	defer c.Stop()
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.SaveOrder(ctx, makeOrder(fmt.Sprintf("%d", i)))
	}
}

func BenchmarkCache_GetOrder_Hit(b *testing.B) {
	c := NewCache(5*time.Minute, 1_000_000)
	defer c.Stop()
	ctx := context.Background()
	// preload
	const preload = 100_000
	ids := make([]string, preload)
	for i := 0; i < preload; i++ {
		ids[i] = fmt.Sprintf("%d", i)
		_ = c.SaveOrder(ctx, makeOrder(ids[i]))
	}
	b.ReportAllocs()
	b.ResetTimer()
	idx := 0
	for i := 0; i < b.N; i++ {
		id := ids[idx]
		if _, err := c.GetOrder(ctx, id); err != nil {
			b.Fatalf("unexpected miss: %v", err)
		}
		idx++
		if idx == preload {
			idx = 0
		}
	}
}

func BenchmarkCache_GetOrder_Parallel(b *testing.B) {
	c := NewCache(5*time.Minute, 1_000_000)
	defer c.Stop()
	ctx := context.Background()
	// preload
	const preload = 100_000
	ids := make([]string, preload)
	for i := 0; i < preload; i++ {
		ids[i] = fmt.Sprintf("%d", i)
		_ = c.SaveOrder(ctx, makeOrder(ids[i]))
	}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			id := ids[r.Intn(preload)]
			if _, err := c.GetOrder(ctx, id); err != nil {
				b.Fatalf("unexpected miss: %v", err)
			}
		}
	})
}
