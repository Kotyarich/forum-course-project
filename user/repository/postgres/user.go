package postgres

import (
	"context"
	"dbProject/db"
	"dbProject/models"
	userPkg "dbProject/user"
	"github.com/jackc/pgx"
)

type User struct {
	About    string
	Email    string
	Fullname string
	Nickname string
}

type UserRepository struct {
	db *pgx.ConnPool
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db.GetDB()}
}

func (r UserRepository) CreateUser(ctx context.Context, user *models.User) ([]*models.User, error) {
	model := toPostgresUser(user)

	_, err := r.db.Exec("INSERT INTO users (about, email, fullname, nickname) "+
		"VALUES ($1, $2, $3, $4)",
		model.About, model.Email, model.Fullname, model.Nickname)
	if err != nil {
		rows, err := r.db.Query("SELECT about, email, fullname, nickname "+
			"FROM users WHERE nickname = $1 OR email = $2",
			model.Nickname, model.Email)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var conflicts []*models.User
		for rows.Next() {
			var u User
			_ = rows.Scan(&u.About, &u.Email, &u.Fullname, &u.Nickname)
			conflicts = append(conflicts, ToModel(&u))
		}

		return conflicts, userPkg.ErrUserAlreadyExists
	}

	return nil, nil
}

func (r UserRepository) AuthUser(ctx context.Context, username, password string) (*models.User, error) {
	// TODO add implementation
	return &models.User{}, nil
}

func (r UserRepository) GetUser(ctx context.Context, username string) (*models.User, error) {
	row := r.db.QueryRow("SELECT about, email, fullname, nickname "+
		"FROM users WHERE nickname = $1", username)

	var user User
	err := row.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)

	if err != nil {
		return nil, userPkg.ErrUserNotFound
	}

	return ToModel(&user), nil
}

func (r UserRepository) ChangeUser(ctx context.Context, user *models.User) (*models.User, error) {
	model := toPostgresUser(user)

	result, err := r.db.Exec("UPDATE users "+
		"SET about = $1, email = $2, fullname = $3 "+
		"WHERE  nickname = $4",
		model.About, model.Email, model.Fullname, model.Nickname)

	if err != nil {
		return nil, userPkg.ErrUserAlreadyExists
	}

	number := result.RowsAffected()
	if number == 0 {
		return nil, userPkg.ErrUserNotFound
	}

	return user, nil
}

func toPostgresUser(u *models.User) *User {
	return &User{
		Nickname: u.Nickname,
		Fullname: u.Fullname,
		Email:    u.Email,
		About:    u.About,
	}
}

func ToModel(u *User) *models.User {
	return &models.User{
		Nickname: u.Nickname,
		Fullname: u.Fullname,
		Email:    u.Email,
		About:    u.About,
	}
}
