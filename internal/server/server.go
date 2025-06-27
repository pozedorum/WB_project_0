package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"wb_project_0/internal/cache"
	"wb_project_0/internal/db"
	"wb_project_0/internal/frontend"
	"wb_project_0/internal/models"

	"github.com/segmentio/kafka-go"
)

type Server struct {
	cache    *cache.Cache
	reader   *kafka.Reader
	server   *http.Server
	frontend *frontend.Frontend

	logger  *log.Logger
	logFile *os.File
}

func New(db *db.Database, kafkaBrokers []string, useKafkaStub bool) *Server {

	s := &Server{}
	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil
	}
	s.logFile = logFile
	s.logger = log.New(os.Stdout, "[SRV] ", log.LstdFlags)
	s.logger.Printf("logger started")
	frnt, err := frontend.New(s.logger)
	if err != nil {
		s.logger.Printf("error with starting frontend")
		return nil
	}
	s.frontend = frnt
	s.logger.Printf("frontend started")
	cach, err := cache.New(db)
	if err != nil {
		s.logger.Printf("error with starting cache")
		return nil
	}
	s.cache = cach
	s.logger.Printf("cache started")
	if useKafkaStub {
		go s.kafkaStub()
		s.logger.Printf("kafkaStub started")
	} else if len(kafkaBrokers) > 0 {
		s.reader = kafka.NewReader(kafka.ReaderConfig{
			Brokers: kafkaBrokers,
			Topic:   "orders",
			GroupID: "order-service",
		})
		go s.consumeKafkaMessages()
		s.logger.Printf("kafka started")
	}

	return s
}

func (s *Server) consumeKafkaMessages() {
	for {
		msg, err := s.reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Kafka read error: %v", err)
			continue
		}

		order, err := s.parseOrder(msg.Value)
		if err != nil {
			log.Printf("Parse error: %v", err)
			continue
		}

		s.cache.Insert(*order)
	}
}

func (s *Server) parseOrder(data []byte) (*models.Order, error) {
	s.logger.Printf("start parsing order")
	var order models.Order
	if err := json.Unmarshal(data, &order); err != nil {
		s.logger.Printf("error with parsing order")
		return nil, err
	}
	return &order, nil
}

func (s *Server) Run(addr string) error {
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.setupRoutes(), // Используем единый роутер
	}

	log.Printf("Starting server on %s", addr)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.reader != nil {
		if err := s.reader.Close(); err != nil {
			log.Printf("Error closing Kafka reader: %v", err)
		}
	}
	return s.server.Shutdown(ctx)
}

// Kafka заглушка - генерирует тестовые данные
func (s *Server) kafkaStub() {
	log.Println("Using Kafka STUB - generating test data")

	// Генерируем тестовый заказ
	testOrder := models.Order{
		OrderUID:    "test124",
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
			Transaction:  "b563feb7b2b84b6test",
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
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmID:              99,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}

	// Преобразуем в JSON
	_, err := json.Marshal(testOrder)
	if err != nil {
		log.Printf("Failed to marshal test order: %v", err)
		return
	}

	// Имитируем получение сообщения из Kafka
	s.cache.Insert(testOrder)
	log.Printf("Inserted test order with ID: %s", testOrder.OrderUID)

	// // Периодически добавляем тестовые данные (по желанию)
	// ticker := time.NewTicker(30 * time.Second)
	// defer ticker.Stop()

	// for range ticker.C {
	// 	// Можно генерировать новые тестовые заказы
	// }
}
