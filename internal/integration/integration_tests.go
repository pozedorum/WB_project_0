package integration

// import (
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"
// 	"wb_project_0/internal/cache"
// 	"wb_project_0/internal/db"
// 	"wb_project_0/internal/models"
// 	"wb_project_0/internal/server"

// 	"github.com/segmentio/kafka-go"
// 	"github.com/stretchr/testify/require"
// )

// func TestOrderFlow(t *testing.T) {
// 	// 1. Инициализируем тестовую БД
// 	testDB, err := db.InitDB()
// 	require.NoError(t, err)
// 	defer testDB.Close()

// 	// 2. Создаем тестовый Kafka Producer
// 	kafkaWriter := &kafka.Writer{
// 		Addr:     kafka.TCP("localhost:9092"),
// 		Topic:    "orders",
// 		Balancer: &kafka.LeastBytes{},
// 	}
// 	defer kafkaWriter.Close()

// 	// 3. Запускаем сервер (в тестовом режиме)
// 	testCache := cache.New(testDB)
// 	srv := server.New(testDB, []string{"localhost:9092"}, "./internal/frontend/templates", "./static")
// 	testServer := httptest.NewServer(srv.setupRoutes())
// 	defer testServer.Close()

// 	// 4. Подготавливаем тестовый заказ
// 	testOrder := models.Order{
// 		OrderUID:    "test123",
// 		TrackNumber: "WBILMTESTTRACK",
// 		Delivery: models.Delivery{
// 			Name:    "Test Testov",
// 			Phone:   "+9720000000",
// 			Zip:     "2639809",
// 			City:    "Kiryat Mozkin",
// 			Address: "Ploshad Mira 15",
// 			Region:  "Kraiot",
// 			Email:   "test@gmail.com",
// 		},
// 		Payment: models.Payment{
// 			Transaction:  "test123",
// 			Currency:     "USD",
// 			Provider:     "wbpay",
// 			Amount:       1817,
// 			PaymentDt:    1637907727,
// 			Bank:         "alpha",
// 			DeliveryCost: 1500,
// 			GoodsTotal:   317,
// 		},
// 		Items: []models.Item{
// 			{
// 				ChrtID:      9934930,
// 				TrackNumber: "WBILMTESTTRACK",
// 				Price:       453,
// 				Rid:         "ab4219087a764ae0btest",
// 				Name:        "Mascaras",
// 				Sale:        30,
// 				Size:        "0",
// 				TotalPrice:  317,
// 			},
// 		},
// 	}

// 	// 5. Отправляем заказ в Kafka
// 	orderJSON, err := json.Marshal(testOrder)
// 	require.NoError(t, err)

// 	err = kafkaWriter.WriteMessages(context.Background(),
// 		kafka.Message{
// 			Key:   []byte(testOrder.OrderUID),
// 			Value: orderJSON,
// 		},
// 	)
// 	require.NoError(t, err)

// 	// 6. Даем время на обработку (можно заменить на более надежный механизм)
// 	time.Sleep(2 * time.Second)

// 	// 7. Проверяем, что заказ появился в БД
// 	dbOrder, err := testDB.GetOrderByUID(context.Background(), testOrder.OrderUID)
// 	require.NoError(t, err)
// 	require.Equal(t, testOrder.OrderUID, dbOrder.OrderUID)

// 	// 8. Проверяем API endpoint
// 	resp, err := http.Get(testServer.URL + "/api/order/" + testOrder.OrderUID)
// 	require.NoError(t, err)
// 	defer resp.Body.Close()
// 	require.Equal(t, http.StatusOK, resp.StatusCode)

// 	var apiOrder models.Order
// 	err = json.NewDecoder(resp.Body).Decode(&apiOrder)
// 	require.NoError(t, err)
// 	require.Equal(t, testOrder.OrderUID, apiOrder.OrderUID)
// }
