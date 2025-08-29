package examples

import "stresser/internal/models"

// Orders contains valid-looking orders that should pass validation.
var Orders = []models.Order{
	{
		OrderUID:    "550e8400-e29b-41d4-a716-446655440000",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "John Doe",
			Phone:   "+12345678901",
			Zip:     "123456",
			City:    "New York",
			Address: "1st Avenue, 10",
			Region:  "NY",
			Email:   "john.doe@example.com",
		},
		Payment: models.Payment{
			Transaction:  "txn0001",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1000,
			PaymentDT:    1637910139,
			Bank:         "Sber",
			DeliveryCost: 150,
			GoodsTotal:   1000,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       100,
				RID:         "ab4219087a764ae0btest",
				Name:        "Medical Mask",
				Sale:        10,
				Size:        "0",
				TotalPrice:  90,
				NmID:        2389212,
				Brand:       "WB",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "customer1",
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmID:              99,
		DateCreated:       "2021-11-26T06:22:19Z",
		OofShard:          "1",
	},
	{
		OrderUID:    "3fa85f64-5717-4562-b3fc-2c963f66afa6",
		TrackNumber: "TRACK0002",
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "Jane Smith",
			Phone:   "+447911123456",
			Zip:     "90210",
			City:    "Los Angeles",
			Address: "Sunset Blvd 42",
			Region:  "CA",
			Email:   "jane.smith@example.com",
		},
		Payment: models.Payment{
			Transaction:  "txn0002",
			Currency:     "EUR",
			Provider:     "stripe",
			Amount:       2500,
			PaymentDT:    1638006539,
			Bank:         "Chase",
			DeliveryCost: 200,
			GoodsTotal:   2500,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:      1234567,
				TrackNumber: "TRACK0002",
				Price:       2500,
				RID:         "rid-0002",
				Name:        "Sneakers",
				Sale:        0,
				Size:        "42",
				TotalPrice:  2500,
				NmID:        7654321,
				Brand:       "Nike",
				Status:      100,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "customer2",
		DeliveryService:   "dhl",
		ShardKey:          "3",
		SmID:              77,
		DateCreated:       "2022-01-01T10:00:00Z",
		OofShard:          "2",
	},
}

// OrdersMissingRequired contains orders that miss required fields.
var OrdersMissingRequired = []models.Order{
	{
		OrderUID:    "a1b2c3d4e5f6g7h8i9j0",
		TrackNumber: "TRACK-MISS-ITEMS",
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "Alice",
			Phone:   "+12025550123",
			Zip:     "10001",
			City:    "New York",
			Address: "5th Ave 1",
			Region:  "NY",
			Email:   "alice@example.com",
		},
		Payment: models.Payment{
			Transaction:  "txn-miss-items",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       100,
			PaymentDT:    1637910139,
			Bank:         "Sber",
			DeliveryCost: 10,
			GoodsTotal:   100,
			CustomFee:    0,
		},
		Items:             []models.Item{}, // invalid: min=1
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "customer3",
		DeliveryService:   "meest",
		ShardKey:          "4",
		SmID:              55,
		DateCreated:       "2021-12-01T12:00:00Z",
		OofShard:          "3",
	},
	{
		OrderUID:    "ABCDEF1234567890",
		TrackNumber: "TRACK-MISS-DELIVERY",
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "", // invalid: required
			Phone:   "+12025550124",
			Zip:     "10002",
			City:    "New York",
			Address: "5th Ave 2",
			Region:  "NY",
			Email:   "bob@example.com",
		},
		Payment: models.Payment{
			Transaction:  "txn-miss-delivery",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       200,
			PaymentDT:    1637911139,
			Bank:         "Sber",
			DeliveryCost: 15,
			GoodsTotal:   200,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:      1,
				TrackNumber: "TRACK-MISS-DELIVERY",
				Price:       200,
				RID:         "rid-miss-delivery",
				Name:        "T-Shirt",
				Sale:        0,
				Size:        "M",
				TotalPrice:  200,
				NmID:        123,
				Brand:       "BrandX",
				Status:      10,
			},
		},
		Locale:          "en",
		CustomerID:      "", // invalid: required
		DeliveryService: "meest",
		ShardKey:        "5",
		SmID:            44,
		DateCreated:     "2021-12-02T12:00:00Z",
		OofShard:        "4",
	},
}

// OrdersBadFormats contains orders with format issues (email, phone, currency, date, etc.).
var OrdersBadFormats = []models.Order{
	{
		OrderUID:    "uid_1", // invalid for alphanum and not a uuid4
		TrackNumber: "TRACK-BAD-FMT-1",
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "Charlie",
			Phone:   "12345", // invalid: not e164
			Zip:     "12A45", // invalid: non-numeric
			City:    "Boston",
			Address: "Some St 3",
			Region:  "MA",
			Email:   "not-an-email", // invalid email
		},
		Payment: models.Payment{
			Transaction:  "txn-bad-fmt-1",
			Currency:     "us", // invalid: len!=3 and not uppercase
			Provider:     "wbpay",
			Amount:       300,
			PaymentDT:    1637912139,
			Bank:         "BOA",
			DeliveryCost: 20,
			GoodsTotal:   300,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:      2,
				TrackNumber: "TRACK-BAD-FMT-1",
				Price:       300,
				RID:         "rid-bad-fmt-1",
				Name:        "Hat",
				Sale:        0,
				Size:        "L",
				TotalPrice:  300,
				NmID:        456,
				Brand:       "BrandY",
				Status:      10,
			},
		},
		Locale:          "en-US", // invalid: not alpha
		CustomerID:      "customer4",
		DeliveryService: "fedex",
		ShardKey:        "6",
		SmID:            33,
		DateCreated:     "2021/12/03 12:00:00", // invalid: wrong datetime format
		OofShard:        "5",
	},
}

// OrdersNegativeValues contains orders with negative numbers where gte=0 is required.
var OrdersNegativeValues = []models.Order{
	{
		OrderUID:    "ABC123",
		TrackNumber: "TRACK-NEG-1",
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "Derek",
			Phone:   "+12025550125",
			Zip:     "10003",
			City:    "Seattle",
			Address: "Pine St 10",
			Region:  "WA",
			Email:   "derek@example.com",
		},
		Payment: models.Payment{
			Transaction:  "txn-neg-1",
			Currency:     "USD",
			Provider:     "paypal",
			Amount:       -1, // invalid
			PaymentDT:    1637913139,
			Bank:         "CITI",
			DeliveryCost: -5,   // invalid
			GoodsTotal:   -100, // invalid
			CustomFee:    -2,   // invalid
		},
		Items: []models.Item{
			{
				ChrtID:      3,
				TrackNumber: "TRACK-NEG-1",
				Price:       -10, // invalid
				RID:         "rid-neg-1",
				Name:        "Gloves",
				Sale:        -1, // invalid
				Size:        "S",
				TotalPrice:  -1, // invalid
				NmID:        789,
				Brand:       "BrandZ",
				Status:      10,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "customer5",
		DeliveryService:   "ups",
		ShardKey:          "7",
		SmID:              22,
		DateCreated:       "2021-12-04T12:00:00Z",
		OofShard:          "6",
	},
}
