package repositories

import (
	"context"
	"database/sql"
	"log"
	"music-service/models"
)

type userRepository struct {
	db *sql.DB
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	// если роль пустая, подставляем "user"
	if user.Role == "" {
		user.Role = "user"
	}

	query := `INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id`

	// Выполняем вставку и сразу сканируем ID
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password, user.Role).Scan(&user.ID)
	if err != nil {
		// Временно логируем реальную ошибку
		log.Println("REAL DB ERROR:", err)
		return err
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, role FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
	return user, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, role FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
	return user, err
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, role FROM users WHERE username = $1`
	err := r.db.QueryRowContext(ctx, query, username).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
	return user, err
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `UPDATE users SET username = $1, email = $2, password = $3, role = $4 WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.Password, user.Role, user.ID)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]models.User, error) {
	query := `SELECT id, username, email, role FROM users LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Role); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
