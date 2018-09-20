package postgres

import (
	"context"
	"errors"

	"github.com/Iliad/vacancytest/pkg/models"
)

func (pgdb *pgDB) GetVacancy(ctx context.Context, id int) (*models.Vacancy, error) {
	pgdb.log.WithField("id", id).Infoln("Getting vacancy")
	var vacancy models.Vacancy
	rows, err := pgdb.conn.QueryxContext(ctx, "SELECT * FROM vacancies WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		if rows.Err() != nil {
			return nil, rows.Err()
		} else {
			return nil, errors.New("vacancy not exists")
		}
	}
	err = rows.StructScan(&vacancy)
	return &vacancy, err
}

func (pgdb *pgDB) GetVacancies(ctx context.Context) ([]models.Vacancy, error) {
	pgdb.log.Infoln("Getting vacancies")
	var vacancies []models.Vacancy
	rows, err := pgdb.conn.QueryxContext(ctx, "SELECT * FROM vacancies ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var vacancy models.Vacancy
		if err := rows.StructScan(&vacancy); err != nil {
			return nil, err
		}
		vacancies = append(vacancies, vacancy)
	}
	return vacancies, err
}

func (pgdb *pgDB) AddVacancy(ctx context.Context, vacancy *models.Vacancy) error {
	pgdb.log.WithField("name", vacancy.Name).Infoln("Creating vacancy")
	rows, err := pgdb.conn.QueryxContext(ctx, "INSERT INTO vacancies (name, salary, experience, city) "+
		"VALUES ($1, $2, $3, $4) RETURNING id",
		vacancy.Name, vacancy.Salary, vacancy.Experience, vacancy.City)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return rows.Err()
	}
	return rows.Scan(&vacancy.ID)
}

func (pgdb *pgDB) DeleteVacancy(ctx context.Context, id int) error {
	pgdb.log.WithField("id", id).Infoln("Getting vacancy")
	res, err := pgdb.conn.ExecContext(ctx, "DELETE FROM vacancies WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rows == 0 {
		return errors.New("vacancy not exists")
	}

	return nil
}
