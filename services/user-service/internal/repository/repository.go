package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type IRepository interface {
	Create(username string, password string) (int64, error)
	GetAll() ([]User, error)
	GetById(id int64) (User, error)
	AuthenticateUser(username, password string) (*User, error)
	VerifyToken(ctx context.Context, userID string, token string) (bool, error)
	CreateToken(userID int64, token string) error
	DeleteToken(token string) error
}

type Repository struct {
	Db *sqlx.DB
}

func NewRepository() *Repository {
	db := newDb()

	createUserTable(db)
	createTokenTable(db)

	r := Repository{Db: db}
	r.Create("user", "user")

	return &r
}

func (r *Repository) Create(username string, password string) (int64, error) {
	passwordHash, err := HashPassword(password)
	if err != nil {
		return 0, err
	}

	query := (`INSERT INTO "user" ("username", "password_hash") VALUES (@username, @password_hash) RETURNING "id";`)
	res, err2 := r.Db.Exec(query, sql.Named("username", username), sql.Named("password_hash", passwordHash))
	if err2 != nil {
		return 0, err2
	}

	return res.LastInsertId()
}

func (r *Repository) GetAll() ([]User, error) {
	users := []User{}
	err := r.Db.Select(&users, `SELECT * FROM "user"`)
	return users, err
}

func (r *Repository) GetById(id int64) (User, error) {
	user := User{}
	err := r.Db.Get(&user, `SELECT * FROM "user" WHERE "id" = $1`, id)
	return user, err
}

func (r *Repository) AuthenticateUser(username, password string) (*User, error) {
	var user User
	err := r.Db.Get(&user, `SELECT * FROM "user" WHERE "username" = ?`, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if !CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

func createUserTable(db *sqlx.DB) {
	query := `CREATE TABLE IF NOT EXISTS "user" (
		"id" INTEGER PRIMARY KEY, 
		"username" VARCHAR(64) NOT NULL, 
		"password_hash" VARCHAR(64) NOT NULL,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP); 
		
		CREATE UNIQUE INDEX IF NOT EXISTS "user_username_idx" ON user (username);`

	_, err := db.Exec(query)
	if err != nil {
		panic("Could not create user table: " + err.Error())
	}
}

func createTokenTable(db *sqlx.DB) {
	query := `CREATE TABLE IF NOT EXISTS "token" (
		"id" INTEGER PRIMARY KEY,
		"user_id" INTEGER NOT NULL,
		"token" VARCHAR(255) NOT NULL,
		"expires_at" TIMESTAMP NOT NULL,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
	);
	CREATE INDEX IF NOT EXISTS "token_user_id_idx" ON token (user_id);
	CREATE INDEX IF NOT EXISTS "token_token_idx" ON token (token);
	CREATE INDEX IF NOT EXISTS "token_expires_at_idx" ON token (expires_at);`

	_, err := db.Exec(query)
	if err != nil {
		panic("Could not create token table: " + err.Error())
	}
}

func (r *Repository) VerifyToken(ctx context.Context, userID string, token string) (bool, error) {
	var count int
	err := r.Db.GetContext(ctx, &count, `
		SELECT COUNT(*) FROM token 
		WHERE user_id = ? AND token = ? AND expires_at > ?`,
		userID, token, time.Now())

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *Repository) CreateToken(userID int64, token string) error {
	expiresAt := time.Now().Add(time.Hour)
	_, err := r.Db.Exec(`
		INSERT INTO token (user_id, token, expires_at)
		VALUES (?, ?, ?)`,
		userID, token, expiresAt)
	return err
}

func (r *Repository) DeleteToken(token string) error {
	_, err := r.Db.Exec(`
		DELETE FROM token
		WHERE token = ?`,
		token)
	return err
}

func (r *Repository) RefreshToken(token string) error {
	expiresAt := time.Now().Add(time.Hour)
	_, err := r.Db.Exec(`
		UPDATE token 
		SET expires_at = ?
		WHERE token = ?`,
		expiresAt, token)
	return err
}
