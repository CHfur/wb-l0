package pgsql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"order-service/config"
	"order-service/internal/domain/models"
	"order-service/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(conf *config.DBConfig) (*Storage, error) {
	const op = "storage.pgsql.New"

	db, err := sql.Open("pgx", fmt.Sprintf("postgres://%s:%s@%s:%v/%s", conf.User, conf.Password, conf.Host, conf.Port, conf.Name))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Orders(ctx context.Context) (map[string]*models.Order, error) {
	const op = "storage.pgsql.Orders"

	stmt, err := s.db.Prepare("SELECT data FROM orders")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, _ := stmt.QueryContext(ctx)

	orders := make(map[string]*models.Order)

	for rows.Next() {
		var order []byte
		err = rows.Scan(&order)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		var orderModel *models.Order

		json.Unmarshal(order, &orderModel)

		orders[orderModel.OrderUid] = orderModel
	}

	return orders, nil
}

func (s *Storage) SaveOrder(ctx context.Context, orderData *models.Order) (err error) {
	const op = "storage.pgsql.SaveOrder"

	stmt, err := s.db.Prepare("INSERT INTO orders (uid, data) VALUES($1, $2)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, orderData.OrderUid, orderData)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrOrderExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
