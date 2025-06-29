package kafkaproducer

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"
	"wb_project_0/internal/models"

	"github.com/segmentio/kafka-go"
)

func main() {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    "orders",
		Balancer: &kafka.Hash{},
	}

	for {
		order := generateOrder() // Используйте вашу функцию генерации
		msg, _ := json.Marshal(order)

		err := writer.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(order.OrderUID),
				Value: msg,
			},
		)

		if err != nil {
			log.Printf("Failed to send message: %v", err)
		} else {
			log.Printf("Sent order: %s", order.OrderUID)
		}

		time.Sleep(10 * time.Second) // Интервал отправки
	}
}

func generateOrder() *models.Order {

	order := getSampleOrders(time.Now().Nanosecond())
	order.OrderUID += strconv.Itoa(time.Now().Nanosecond() % 100)
	return order
}

func getSampleOrders(ind int) *models.Order {
	orderSamples := []*models.Order{
		// 1. Стандартный заказ
		{
			OrderUID:    "test1",
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
				Transaction:  "test1",
				Currency:     "USD",
				Provider:     "wbpay",
				Amount:       1817,
				PaymentDt:    1637907727,
				Bank:         "alpha",
				DeliveryCost: 1500,
				GoodsTotal:   317,
				CustomFee:    0,
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
			DateCreated: time.Now(),
		},

		// 2. Заказ с двумя товарами
		{
			OrderUID:    "test2",
			TrackNumber: "WBILMTRACK002",
			Entry:       "WBIL",
			Delivery: models.Delivery{
				Name:    "Ivan Ivanov",
				Phone:   "+79161234567",
				Zip:     "101000",
				City:    "Moscow",
				Address: "Tverskaya st. 10",
				Region:  "Moscow",
				Email:   "ivanov@mail.ru",
			},
			Payment: models.Payment{
				Transaction:  "test2",
				Currency:     "RUB",
				Provider:     "wbpay",
				Amount:       4500,
				PaymentDt:    1637910000,
				Bank:         "sber",
				DeliveryCost: 500,
				GoodsTotal:   4000,
				CustomFee:    0,
			},
			Items: []models.Item{
				{
					ChrtID:      8823411,
					TrackNumber: "WBILMTRACK002",
					Price:       2500,
					Rid:         "bc5320198b875bf1demo",
					Name:        "Smartphone",
					Sale:        10,
					Size:        "0",
					TotalPrice:  2250,
					NmID:        3344556,
					Brand:       "Xiaomi",
					Status:      202,
				},
				{
					ChrtID:      9923412,
					TrackNumber: "WBILMTRACK002",
					Price:       2000,
					Rid:         "de5320198c875bf2demo",
					Name:        "Headphones",
					Sale:        20,
					Size:        "0",
					TotalPrice:  1600,
					NmID:        4455667,
					Brand:       "Sony",
					Status:      202,
				},
			},
			DateCreated: time.Now().Add(-24 * time.Hour),
		},
		// 3. Стандарный заказ №2
		{
			OrderUID:    "test3",
			TrackNumber: "WBILMTESTTRACK",
			Entry:       "WBIL",
			Delivery: models.Delivery{
				Name:    "Gol Golovich",
				Phone:   "+9720000000",
				Zip:     "2639809",
				City:    "Moskow",
				Address: "Ploshad Goidu 21",
				Region:  "Moskow",
				Email:   "test@gmail.com",
			},
			Payment: models.Payment{
				Transaction:  "test3",
				Currency:     "USD",
				Provider:     "wbpay",
				Amount:       1817,
				PaymentDt:    1637907727,
				Bank:         "alpha",
				DeliveryCost: 1500,
				GoodsTotal:   317,
				CustomFee:    0,
			},
			Items: []models.Item{
				{
					ChrtID:      9934931,
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
			DateCreated: time.Now(),
		},
	}
	return orderSamples[ind%3]
}
