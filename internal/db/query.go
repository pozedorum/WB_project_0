package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"wb_project_0/internal/models"
)

func (db *Database) SaveOrder(ctx context.Context, order models.Order) error {
	//db.logger.Printf("Saving order %s", order.OrderUID)
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		db.logger.Printf("Begin transaction failed: %v", err)
		return fmt.Errorf("begin transaction failed: %w", err)
	}
	defer func() {
		if err != nil {
			db.logger.Printf("Rolling back transaction: %v", err)
			tx.Rollback()
		}
	}()
	if err = CheckFieldsOrder(order); err != nil {
		//db.logger.Printf("The order is invalid: %v", err)
		return err
	}

	if err != nil {
		db.logger.Printf("%v", err) // Non-valid structure
	}
	// Сохраняем основные данные заказа
	//db.logger.Printf("Inserting order %s", order.OrderUID)
	if err := insertOrder(ctx, tx, order); err != nil {
		db.logger.Printf("Inserting order failed %s", err)
		return err
	}

	// Сохраняем доставку
	//db.logger.Printf("Inserting delivery %s", order.OrderUID)
	if err := insertDelivery(ctx, tx, order); err != nil {
		db.logger.Printf("Inserting delivery failed %s", err)
		return err
	}

	// Сохраняем платежи
	//db.logger.Printf("Inserting payment %s", order.OrderUID)
	if err := insertPayment(ctx, tx, order); err != nil {
		db.logger.Printf("Inserting payment failed %s", err)
		return err
	}

	// Сохраняем товары
	//db.logger.Printf("Inserting items %s", order.OrderUID)
	if err := insertItems(ctx, tx, order); err != nil {
		db.logger.Printf("Inserting items failed %s", err)
		return err
	}

	return tx.Commit()
}

func insertOrder(ctx context.Context, tx *sql.Tx, order models.Order) error {
	query := `
	INSERT INTO orders (
		order_uid, track_number, entry, locale, 
		internal_signature, customer_id, delivery_service,
		shardkey, sm_id, date_created, oof_shard
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	ON CONFLICT (order_uid) DO NOTHING`

	_, err := tx.ExecContext(ctx, query,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	)

	return err
}

func insertDelivery(ctx context.Context, tx *sql.Tx, order models.Order) error {
	query := `
	INSERT INTO deliveries (
		order_uid, name, phone, zip, 
		city, address, region,
		email
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT (order_uid) DO UPDATE SET
		name = EXCLUDED.name,
		phone = EXCLUDED.phone,
		zip = EXCLUDED.zip,
		city = EXCLUDED.city,
		address = EXCLUDED.address,
		region = EXCLUDED.region,
		email = EXCLUDED.email`

	_, err := tx.ExecContext(ctx, query,
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)

	return err
}

func insertPayment(ctx context.Context, tx *sql.Tx, order models.Order) error {
	query := `
	INSERT INTO payments (
		order_uid, transaction, request_id, currency, provider, 
		amount, payment_dt, bank,
		delivery_cost, goods_total, custom_fee
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	ON CONFLICT (order_uid) DO UPDATE SET
		order_uid = EXCLUDED.order_uid,
		transaction = EXCLUDED.transaction,
		request_id = EXCLUDED.request_id,
		currency = EXCLUDED.currency,
		provider = EXCLUDED.provider,
		amount = EXCLUDED.amount,
		payment_dt = EXCLUDED.payment_dt,
		bank = EXCLUDED.bank,
		delivery_cost = EXCLUDED.delivery_cost,
		goods_total = EXCLUDED.goods_total,
		custom_fee = EXCLUDED.custom_fee`

	_, err := tx.ExecContext(ctx, query,
		order.OrderUID,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)

	return err
}

func insertItems(ctx context.Context, tx *sql.Tx, order models.Order) error {
	var err error
	// Удаляем товары заказа, если они есть
	if _, err = tx.ExecContext(ctx, "DELETE FROM items WHERE order_uid = $1", order.OrderUID); err != nil {
		return err
	}

	// Вставляем новые товары
	query := `
	INSERT INTO items (
		order_uid, chrt_id, track_number, price, rid,
		name, sale, size, total_price, nm_id,
		brand, status
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, query,
			order.OrderUID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)
		if err != nil {
			break
		}
	}
	return err
}

func (db *Database) GetOrderByUID(ctx context.Context, orderUID string) (*models.Order, error) {
	var order models.Order

	//db.logger.Printf("Getting order %s", orderUID)

	err := db.conn.QueryRowContext(ctx, "SELECT * FROM orders WHERE order_uid = $1", orderUID).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
	)
	if err != nil {
		db.logger.Printf("Error 1 getting order %s: %v", orderUID, err)
		return nil, err
	}

	err = db.conn.QueryRowContext(ctx, "SELECT * FROM deliveries WHERE order_uid = $1", orderUID).Scan(
		&order.OrderUID,
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.Zip,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email,
	)
	if err != nil {
		db.logger.Printf("Error 2 getting order %s: %v", orderUID, err)
		return nil, err
	}

	err = db.conn.QueryRowContext(ctx,
		"SELECT * FROM payments WHERE order_uid = $1", orderUID).Scan(
		&order.OrderUID,
		&order.Payment.Transaction,
		&order.Payment.RequestID,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.PaymentDt,
		&order.Payment.Bank,
		&order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
	)
	if err != nil {
		db.logger.Printf("Error 3 getting order %s: %v", orderUID, err)
		return nil, err
	}

	rows, err := db.conn.QueryContext(ctx,
		"SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE order_uid = $1", orderUID)
	if err != nil {
		db.logger.Printf("Error 4 getting order %s: %v", orderUID, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			db.logger.Printf("Error 5 getting order %s: %v", orderUID, err)
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	//db.logger.Printf("Successfully retrieved order %s", orderUID)
	return &order, nil
}

func (db *Database) GetAllOrders(ctx context.Context) ([]models.Order, error) {

	db.logger.Printf("Getting all orders")

	rows, err := db.conn.QueryContext(ctx, "SELECT order_uid FROM orders")
	if err != nil {
		db.logger.Printf("Error getting order_uid: %v", err)
		return nil, err
	}

	defer rows.Close()
	var orders []models.Order
	for rows.Next() {
		var orderUID string
		if err := rows.Scan(&orderUID); err != nil {
			db.logger.Printf("Error orderUID not found: %v", err)
			return nil, err
		}

		order, err := db.GetOrderByUID(ctx, orderUID)
		if err != nil {
			db.logger.Printf("Error not found order with orderUID %s: %v", orderUID, err)
			return nil, err
		}

		orders = append(orders, *order)
	}

	return orders, nil

}

func CheckFieldsOrder(order models.Order) error {
	if order.OrderUID == "" {
		return fmt.Errorf("order UID cannot be empty")
	}
	if order.TrackNumber == "" {
		return fmt.Errorf("track number cannot be empty")
	}
	if len(order.Items) == 0 {
		return fmt.Errorf("order must contain at least one item")
	}

	for _, item := range order.Items {
		if err := CheckFieldsItem(item); err != nil {
			return err
		}
	}
	return nil
}

func CheckFieldsDelivery(obj models.Delivery) error {
	if obj.Name == "" || obj.Phone == "" || obj.Zip == "" || obj.City == "" ||
		obj.Address == "" || obj.Region == "" || obj.Email == "" {
		return errors.New("non-valid structure Delivery")
	}
	return nil
}

func CheckFieldsPayment(obj models.Payment) error {
	if obj.Transaction == "" || obj.RequestID == "" || obj.Currency == "" ||
		obj.Provider == "" || obj.Amount < 0 || obj.Bank == "" ||
		obj.DeliveryCost < 0 || obj.GoodsTotal < 0 || obj.CustomFee < 0 {
		return errors.New("non-valid structure Payment")
	}
	return nil
}

func CheckFieldsItem(obj models.Item) error {
	if obj.ChrtID == 0 {
		return fmt.Errorf("item chrt_id cannot be zero")
	}
	if obj.Price <= 0 {
		return fmt.Errorf("item price must be positive")
	}
	return nil
}
