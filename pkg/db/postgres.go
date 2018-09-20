package db

import (
	"io"

	"github.com/Iliad/vacancytest/pkg/models"

	"context"
)

// DB is an interface for persistent data storage (also sometimes called DAO).
type DB interface {
	GetVacancy(ctx context.Context, id int) (*models.Vacancy, error)
	GetVacancies(ctx context.Context) ([]models.Vacancy, error)
	AddVacancy(ctx context.Context, vacancy *models.Vacancy) error
	DeleteVacancy(ctx context.Context, id int) error

	GetUser(ctx context.Context, login string) (*models.User, error)
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUsersCount(ctx context.Context) (*uint, error)
	CreateUser(ctx context.Context, user *models.User) error
	ChangeUserRole(ctx context.Context, login string, role models.UserRole) error

	io.Closer
}
