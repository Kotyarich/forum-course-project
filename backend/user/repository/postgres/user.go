package postgres

import (
	"context"
	"crypto/sha1"
	"dbProject/db"
	"dbProject/models"
	userPkg "dbProject/user"
	"fmt"
	"github.com/jackc/pgx"
	"strconv"
)

type User struct {
	About    string
	Email    string
	Fullname string
	Password string
	Nickname string
	IsAdmin  bool
}

type UserRepository struct {
	db *pgx.ConnPool
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db.GetDB()}
}

func (r UserRepository) CreateUser(ctx context.Context, user *models.User) ([]*models.User, int, error) {
	model := toPostgresUser(user)
	userId := -1

	err := r.db.QueryRow("INSERT INTO users (about, email, fullname, nickname, password) "+
		"VALUES ($1, $2, $3, $4, $5) RETURNING id",
		model.About, model.Email, model.Fullname, model.Nickname, model.Password).Scan(&userId)
	if err != nil {
		rows, err := r.db.Query("SELECT about, email, fullname, nickname "+
			"FROM users WHERE nickname = $1 OR email = $2",
			model.Nickname, model.Email)
		if err != nil {
			return nil, userId, err
		}
		defer rows.Close()

		var conflicts []*models.User
		for rows.Next() {
			var u User
			_ = rows.Scan(&u.About, &u.Email, &u.Fullname, &u.Nickname)
			conflicts = append(conflicts, ToModel(&u))
		}

		return conflicts, userId, userPkg.ErrUserAlreadyExists
	}

	return nil, userId, nil
}

func (r UserRepository) AuthUser(ctx context.Context, username, password string) (*models.User, int, error) {
	row := r.db.QueryRow("SELECT id, about, email, fullname, nickname "+
		"FROM users WHERE nickname = $1 AND password = $2", username, password)

	userId := -1
	var user User
	err := row.Scan(&userId, &user.About, &user.Email, &user.Fullname, &user.Nickname)
	if err != nil {
		fmt.Println(err.Error())
		return nil, userId, userPkg.ErrUserNotFound
	}

	return ToModel(&user), userId, nil
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
		Password: u.Password,
		Email:    u.Email,
		About:    u.About,
		IsAdmin:  u.IsAdmin,
	}
}

func ToModel(u *User) *models.User {
	return &models.User{
		Nickname: u.Nickname,
		Fullname: u.Fullname,
		Email:    u.Email,
		About:    u.About,
		IsAdmin:  u.IsAdmin,
	}
}

func createToken(userId int, userAgent, authTime string) string {
	pwd := sha1.New()
	pwd.Write([]byte(strconv.Itoa(userId)))
	pwd.Write([]byte(userAgent))
	pwd.Write([]byte(authTime))
	return fmt.Sprintf("%x", pwd.Sum(nil))
}
