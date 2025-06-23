package database

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"myblog/internal/db"
	"myblog/internal/models"
)

func TestMain(m *testing.M) {
	// Подготовка тестовой БД
	code := m.Run()
	os.Exit(code)
}

func TestSaveAndGetOrder(t *testing.T) {
	d := setupTestDB(t)
	ctx := context.Background()

	testOrder := models.Order{
		OrderUID:    "test123",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: models.Payment{
			Transaction:  "test123",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    time.Now().Unix(),
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
		},
		Items: []models.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmID:              99,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}

	// Тест сохранения
	err := d.SaveOrder(ctx, testOrder)
	assert.NoError(t, err)

	// Тест чтения
	saved, err := d.GetOrderByUID(ctx, testOrder.OrderUID)
	assert.NoError(t, err)
	assert.Equal(t, testOrder.OrderUID, saved.OrderUID)

	// Тест обновления
	testOrder.Payment.Amount = 999
	err = d.SaveOrder(ctx, testOrder)
	assert.NoError(t, err)

	updated, err := d.GetOrderByUID(ctx, testOrder.OrderUID)
	assert.NoError(t, err)
	assert.Equal(t, 999, updated.Payment.Amount)
}

func TestGetAllOrders(t *testing.T) {
	d := setupTestDB(t)
	ctx := context.Background()

	orders, err := d.GetAllOrders(ctx)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(orders), 1)
}

func setupTestDB(t *testing.T) *db.Database {
	d, err := db.InitDB()
	if err != nil {
		t.Fatalf("Failed to init DB: %v", err)
	}
	t.Cleanup(func() {
		if err := d.Close(); err != nil {
			t.Errorf("Failed to close DB: %v", err)
		}
	})
	return d
}