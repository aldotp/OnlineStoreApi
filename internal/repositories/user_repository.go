package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/aldotp/OnlineStore/internal/entity"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u UserRepository) GetUserByID(id int) (*entity.User, error) {
	row := u.db.QueryRow("SELECT id, username, password, email, created_at, updated_at FROM users WHERE id = ?", id)
	var user entity.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UserRepository) CreateUser(ctx context.Context, cust entity.User) (*entity.User, error) {

	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO users (username, password, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)", cust.Username, cust.Password, cust.Email, time.Now().UTC(), time.Now().UTC())
	if err != nil {
		return nil, err
	}

	row := tx.QueryRowContext(ctx, "SELECT id, username, password, email, created_at, updated_at FROM users WHERE username = ?", cust.Username)
	insertedUser := new(entity.User)
	if err := row.Scan(&insertedUser.ID, &insertedUser.Username, &insertedUser.Password, &insertedUser.Email, &insertedUser.CreatedAt, &insertedUser.UpdatedAt); err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO carts (user_id, created_at, updated_at) VALUES (?, ?, ?)",
		insertedUser.ID,
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		return nil, err
	}

	return insertedUser, tx.Commit()
}

func (u UserRepository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	row := u.db.QueryRowContext(ctx, "SELECT id, username, password, email, created_at, updated_at FROM users WHERE username = ?", username)
	var user entity.User

	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
