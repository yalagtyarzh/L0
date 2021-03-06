package dbrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/yalagtyarzh/L0/internal/models"
)

// InsertOrder inserts order object into database
func (m *postgresDBRepo) InsertOrder(o models.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
			insert into orders (order_uid, track_number, entry, locale, internal_signature, customer_id,
			                    delivery_service, shardkey, sm_id, date_created, oof_shard)
			values
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning id
	`

	var orderID int
	err := m.DB.QueryRowContext(
		ctx,
		query,
		o.OrderUID,
		o.TrackNumber,
		o.Entry,
		o.Locale,
		o.InternalSignature,
		o.CustomerID,
		o.DeliveryService,
		o.Shardkey,
		o.SmID,
		o.DateCreated,
		o.OofShard,
	).Scan(&orderID)

	if err != nil {
		return err
	}

	err = m.InsertPaymentByOrderID(o.Payment, orderID)
	if err != nil {
		return err
	}

	err = m.InsertDeliveryByOrderID(o.Delivery, orderID)
	if err != nil {
		return err
	}

	for _, item := range o.Items {
		id, err := m.InsertItem(item)
		if err != nil {
			return err
		}

		err = m.InsertOrderItems(orderID, id)
		if err != nil {
			return err
		}
	}

	return nil
}

// InsertPaymentByOrderID inserts payment object into database by order's ID
func (m *postgresDBRepo) InsertPaymentByOrderID(pm models.Payment, orderID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
			insert into payment (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost,
			                     goods_total, custom_fee, order_id)
			values
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := m.DB.ExecContext(
		ctx,
		query,
		pm.Transaction,
		pm.RequestID,
		pm.Currency,
		pm.Provider,
		pm.Amount,
		pm.PaymentDt,
		pm.Bank,
		pm.DeliveryCost,
		pm.GoodsTotal,
		pm.CustomFee,
		orderID,
	)

	if err != nil {
		return err
	}

	return nil
}

// InsertDeliveryByOrderID inserts delivery object into database by order's ID
func (m *postgresDBRepo) InsertDeliveryByOrderID(d models.Delivery, orderID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
			insert into delivery (name, phone, zip, city, address, region, email, order_id)
			values
			($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := m.DB.ExecContext(
		ctx,
		query,
		d.Name,
		d.Phone,
		d.Zip,
		d.City,
		d.Address,
		d.Region,
		d.Email,
		orderID,
	)

	if err != nil {
		return err
	}

	return nil
}

// InsertItem inserts item object into database
func (m *postgresDBRepo) InsertItem(i models.Item) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
			insert into items (track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
			values
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id
	`

	var newID int

	err := m.DB.QueryRowContext(
		ctx,
		query,
		i.TrackNumber,
		i.Price,
		i.Rid,
		i.Name,
		i.Sale,
		i.Size,
		i.TotalPrice,
		i.NmID,
		i.Brand,
		i.Status,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertOrderItems connects order with items on it
func (m *postgresDBRepo) InsertOrderItems(orderID, chrtID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
			insert into order_items (order_id, chrt_id)
			values
			($1, $2) 
	`

	_, err := m.DB.ExecContext(ctx, query, orderID, chrtID)
	if err != nil {
		return err
	}

	return nil
}

// GetOrders returns all orders
func (m *postgresDBRepo) GetOrders() ([]models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var orders []models.Order
	query := `
			select id, order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, 
			       sm_id, date_created, oof_shard from orders
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return orders, err
	}

	for rows.Next() {
		var id int
		var order models.Order
		err := rows.Scan(
			&id,
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.Shardkey,
			&order.SmID,
			&order.DateCreated,
			&order.OofShard,
		)
		if err != nil {
			return orders, err
		}

		delivery, err := m.GetDeliveryByOrderUID(id)
		if err != nil {
			return orders, err
		}

		payment, err := m.GetPaymentByOrderUID(id)
		if err != nil {
			return orders, err
		}

		items, err := m.GetItemsByOrderUID(id)
		if err != nil {
			return orders, err
		}

		order.Delivery = delivery
		order.Payment = payment
		order.Items = items

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return orders, err
	}

	return orders, nil
}

// GetDeliveryByOrderUID returns delivery object by order's ID
func (m *postgresDBRepo) GetDeliveryByOrderUID(orderID int) (models.Delivery, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var delivery models.Delivery
	query := `
			select name, phone, zip, city, address, region, email from delivery
			where order_id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, orderID)
	err := row.Scan(
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email,
	)

	if err != nil && err != sql.ErrNoRows {
		return delivery, err
	}

	return delivery, nil
}

// GetPaymentByOrderUID returns payment object by order's ID
func (m *postgresDBRepo) GetPaymentByOrderUID(orderID int) (models.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var payment models.Payment
	query := `
			select "transaction", request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total,
			custom_fee
			from payment
			where order_id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, orderID)
	err := row.Scan(
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDt,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee,
	)

	if err != nil && err != sql.ErrNoRows {
		return payment, err
	}

	return payment, nil
}

// GetItemsByOrderUID returns all items object by order's ID
func (m *postgresDBRepo) GetItemsByOrderUID(orderID int) ([]models.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var items []models.Item
	query := `
			select i.track_number, i.price, i.rid, i.name, i.sale, i.size, i.total_price, i.nm_id, i.brand, i.status
			from items i
			join order_items oi on (oi.chrt_id = i.id)
			where oi.order_id = $1
	`

	rows, err := m.DB.QueryContext(ctx, query, orderID)
	if err != nil {
		return items, err
	}

	for rows.Next() {
		var item models.Item
		err := rows.Scan(
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
			return items, err
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return items, err
	}

	return items, nil
}
