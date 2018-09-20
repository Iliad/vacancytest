package postgres

import (
	"github.com/Iliad/vacancytest/pkg/db"
	"github.com/golang-migrate/migrate"
	migdrv "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file" // needed to load migrations scripts from files
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgresql database driver
	"github.com/sirupsen/logrus"
)

type pgDB struct {
	conn *sqlx.DB // do not use it in select/exec operations
	log  *logrus.Entry
}

// DBConnect initializes connection to postgresql database.
// github.com/jmoiron/sqlx used to to get work with database.
// Function tries to ping database and apply migrations using github.com/mattes/migrate.
// If migrations applying failed database goes to dirty state and requires manual conflict resolution.
func DBConnect(pgConnStr string, migrationsPath string) (db.DB, error) {
	log := logrus.WithField("component", "db")
	log.Infoln("Connecting to ", pgConnStr)
	conn, err := sqlx.Open("postgres", pgConnStr)
	if err != nil {
		log.WithError(err).Errorln("Postgres connection failed")
		return nil, err
	}
	if pingErr := conn.Ping(); pingErr != nil {
		return nil, pingErr
	}

	ret := &pgDB{
		conn: conn,
		log:  log,
	}

	m, err := ret.migrateUp(migrationsPath)
	if err != nil {
		return nil, err
	}
	version, _, _ := m.Version()
	log.WithField("version", version).Infoln("Migrate up")

	return ret, nil
}

func (pgdb *pgDB) migrateUp(path string) (*migrate.Migrate, error) {
	pgdb.log.Infof("Running migrations")
	instance, err := migdrv.WithInstance(pgdb.conn.DB, &migdrv.Config{MigrationsTable: "migrations_um"})
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+path, "postgres", instance)
	if err != nil {
		return nil, err
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}
	return m, nil
}

func (pgdb *pgDB) Close() error {
	return pgdb.conn.Close()
}
