package db

import (
	"context"
	"os"
	"testing"
	"time"

	"wb_project_0/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Подготовка тестовой БД
	code := m.Run()
	os.Exit(code)
}

func TestDatabaseCRUD(t *testing.T) {
	d := setupTestDB(t)
	ctx := context.Background()

	// Тестовые данные
	testOrder := createTestOrder("test123")
	updatedOrder := createTestOrder("test123")
	updatedOrder.Payment.Amount = 999

	t.Run("SaveOrder", func(t *testing.T) {
		err := d.SaveOrder(ctx, testOrder)
		assert.NoError(t, err)
	})

	t.Run("GetOrderByUID", func(t *testing.T) {
		order, err := d.GetOrderByUID(ctx, testOrder.OrderUID)
		require.NoError(t, err)
		assert.Equal(t, testOrder.OrderUID, order.OrderUID)
		assert.Equal(t, testOrder.Payment.Amount, order.Payment.Amount)
	})

	t.Run("UpdateOrder", func(t *testing.T) {
		err := d.SaveOrder(ctx, updatedOrder)
		assert.NoError(t, err)

		order, err := d.GetOrderByUID(ctx, testOrder.OrderUID)
		require.NoError(t, err)
		assert.Equal(t, updatedOrder.Payment.Amount, order.Payment.Amount)
	})

	t.Run("GetAllOrders", func(t *testing.T) {
		orders, err := d.GetAllOrders(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(orders), 1)
		assert.Equal(t, updatedOrder.OrderUID, orders[0].OrderUID)
	})

	t.Run("GetNonExistentOrder", func(t *testing.T) {
		_, err := d.GetOrderByUID(ctx, "nonexistent")
		assert.Error(t, err)
	})
}

func TestDatabaseTablesOperations(t *testing.T) {
	d := setupTestDB(t)

	t.Run("CheckTablesExist", func(t *testing.T) {
		exists, err := d.CheckAllTablesExist()
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("DeleteAndRecreateTables", func(t *testing.T) {
		err := d.DeleteTables()
		require.NoError(t, err)

		exists, err := d.CheckAllTablesExist()
		require.NoError(t, err)
		assert.False(t, exists)

		err = d.CreateTables()
		require.NoError(t, err)

		exists, err = d.CheckAllTablesExist()
		require.NoError(t, err)
		assert.True(t, exists)
	})
}

func TestDatabaseConnection(t *testing.T) {
	t.Run("InitDB", func(t *testing.T) {
		db, err := InitDB()
		require.NoError(t, err)
		require.NotNil(t, db)
		assert.NoError(t, db.Close())
	})

	t.Run("CloseConnection", func(t *testing.T) {
		db, err := InitDB()
		require.NoError(t, err)
		assert.NoError(t, db.Close())
	})
}

func TestDatabaseErrorCases(t *testing.T) {
	d := setupTestDB(t)
	ctx := context.Background()

	t.Run("SaveInvalidOrder", func(t *testing.T) {
		invalidOrder := models.Order{} // Пустой заказ
		err := d.SaveOrder(ctx, invalidOrder)
		assert.Error(t, err)
	})

	t.Run("SaveOrderWithInvalidItems", func(t *testing.T) {
		order := createTestOrder("invalid_items")
		order.Items = []models.Item{{}} // Невалидный товар
		err := d.SaveOrder(ctx, order)
		assert.Error(t, err)
	})
}

func TestDatabaseConcurrency(t *testing.T) {
	d := setupTestDB(t)
	ctx := context.Background()

	t.Run("ConcurrentWrites", func(t *testing.T) {
		const numOrders = 5
		orders := make([]models.Order, numOrders)
		for i := 0; i < numOrders; i++ {
			orders[i] = createTestOrder("concurrent_" + string(rune('a'+i)))
		}

		errs := make(chan error, numOrders)
		for _, order := range orders {
			go func(o models.Order) {
				errs <- d.SaveOrder(ctx, o)
			}(order)
		}

		for i := 0; i < numOrders; i++ {
			assert.NoError(t, <-errs)
		}

		savedOrders, err := d.GetAllOrders(ctx)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(savedOrders), numOrders)
	})
}

func createTestOrder(uid string) models.Order {
	return models.Order{
		OrderUID:    uid,
		TrackNumber: "TRACK_" + uid,
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "Test User",
			Phone:   "+1234567890",
			Zip:     "123456",
			City:    "Test City",
			Address: "Test Address",
			Region:  "Test Region",
			Email:   "test@test.com",
		},
		Payment: models.Payment{
			Transaction:  "trans_" + uid,
			Currency:     "USD",
			Provider:     "test_provider",
			Amount:       1000,
			PaymentDt:    time.Now().Unix(),
			Bank:         "test_bank",
			DeliveryCost: 500,
			GoodsTotal:   500,
		},
		Items: []models.Item{
			{
				ChrtID:      1000,
				TrackNumber: "TRACK_" + uid,
				Price:       100,
				Rid:         "rid_" + uid,
				Name:        "Test Item",
				Sale:        10,
				Size:        "0",
				TotalPrice:  90,
				NmID:        2000,
				Brand:       "Test Brand",
				Status:      1,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test_customer",
		DeliveryService:   "test_service",
		ShardKey:          "1",
		SmID:              1,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}
}

func setupTestDB(t *testing.T) *Database {
	d, err := InitDB()
	require.NoError(t, err)

	// Очищаем лог-файл перед каждым тестом
	if d.logFile != nil {
		_ = d.logFile.Close()
		os.Remove("db.log") // Удаляем файл, чтобы начать "чистый" лог
	}

	require.NoError(t, d.DeleteTables())
	require.NoError(t, d.CreateTables())

	t.Cleanup(func() {
		if err := d.Close(); err != nil {
			t.Logf("Cleanup error: %v", err)
		}
	})
	return d
}
