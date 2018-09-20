package postgres

import (
	"context"
	"errors"

	"github.com/Iliad/vacancytest/pkg/models"
)

func (pgdb *pgDB) CreateUser(ctx context.Context, user *models.User) error {
	pgdb.log.WithField("login", user.Login).WithField("role", user.Role).Infoln("Creating user")
	rows, err := pgdb.conn.QueryxContext(ctx, "INSERT INTO users (login, password, role) "+
		"VALUES ($1, $2, $3) RETURNING id",
		user.Login, user.Password, user.Role)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return rows.Err()
	}
	return rows.Scan(&user.ID)
}

func (pgdb *pgDB) GetUsersCount(ctx context.Context) (*uint, error) {
	pgdb.log.Infoln("Count users")
	var usersCount uint
	rows, err := pgdb.conn.QueryxContext(ctx, "SELECT count(id) FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, rows.Err()
	}
	err = rows.Scan(&usersCount)

	return &usersCount, err
}

func (pgdb *pgDB) GetUsers(ctx context.Context) ([]models.User, error) {
	pgdb.log.Infoln("Get users")
	var users []models.User
	rows, err := pgdb.conn.QueryxContext(ctx, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user models.User
		if err := rows.StructScan(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, err
}

func (pgdb *pgDB) GetUser(ctx context.Context, login string) (*models.User, error) {
	pgdb.log.Infoln("Get user", login)
	var user models.User
	rows, err := pgdb.conn.QueryxContext(ctx, "SELECT * FROM users WHERE login = $1", login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		if rows.Err() != nil {
			return nil, rows.Err()
		} else {
			return nil, errors.New("user not exists")
		}
	}
	err = rows.StructScan(&user)
	return &user, err
}

func (pgdb *pgDB) ChangeUserRole(ctx context.Context, login string, role models.UserRole) error {
	pgdb.log.WithField("login", login).Infoln("Changing user role")
	res, err := pgdb.conn.ExecContext(ctx, "UPDATE users SET role = $1 WHERE login = $2",
		role, login)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rows == 0 {
		return errors.New("user not exists")
	}
	return nil
}
