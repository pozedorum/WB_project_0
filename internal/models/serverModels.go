package models

import (
	"sync"
)

// OrderResponse - структура для HTTP-ответа
type OrderResponse struct {
	Order
}

// Cache - структура для кеша в памяти
type Cache struct {
	Orders map[string]Order
	Sync   sync.RWMutex
}

// KafkaMessage - структура входящего сообщения
type KafkaMessage struct {
	Value     []byte
	Partition int32
	Offset    int64
}
