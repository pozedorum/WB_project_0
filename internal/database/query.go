package database

import (
	"context"
	"database/sql"

	"wb_project_0/internal/models"
)


func (db *Database) SaveOrder(ctx context.Context, order models.Order) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction failed: %w", err)
	}
	defer tx.Rollback()

	// Сохраняем основные данные заказа
	if err := insertOrder(ctx, tx, order); err != nil {
		return err
	}

	// Сохраняем доставку
	if err := insertDelivery(ctx, tx, order); err != nil {
		return err
	}

	// Сохраняем платежи
	if err := insertPayment(ctx, tx, order); err != nil {
		return err
	}

	// Сохраняем товары
	if err := insertItems(ctx, tx, order); err != nil {
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
			item.TrackNumber
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

	err := db.conn.QueryRowContext(ctx, "SELECT * FROM orders WHERE order_uid = $1",orderUID).Scan(
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
		&order.OofShard
	)
	if err != nil {
		return nil, err
	}

	err := db.conn.QueryRowContext(ctx, "SELECT * FROM deliveries WHERE order_uid = $1", orderUID).Scan(
		&order.Delivery.OrderUID,
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.Zip,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email,
	)
	if err != nil {
		return nil, err
	}

	err = db.conn.QueryRowContext(ctx,
		"SELECT * FROM payments WHERE order_uid = $1", orderUID).Scan(
		&order.Payment.OrderUID,
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
		return nil, err
	}

	rows, err := db.conn.QueryContext(ctx,
		"SELECT * FROM items WHERE order_uid = $1", orderUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ID,
			&item.OrderUID,
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
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

}



func (db *Database) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	rows, err := db.conn.QueryContext(ctx,"SELECT order_uid FROM orders")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var orderUid string
		if err := db.conn.QueryContext(&orderUid); err != nil {
			return nil, err
		}

		if order, err := GetOrderByUID(ctx, orderUid); err != nil {
			return nil, err
		}

		orders = append(orders, *order)
	}

	return orders, nil

}