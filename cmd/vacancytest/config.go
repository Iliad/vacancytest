package main

import (
	"fmt"

	"gopkg.in/urfave/cli.v1/altsrc"

	"github.com/Iliad/vacancytest/pkg/db"
	"github.com/Iliad/vacancytest/pkg/db/postgres"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
)

const (
	configFlag       = "config"
	portFlag         = "port"
	debugFlag        = "debug"
	textlogFlag      = "textlog"
	dbPGLoginFlag    = "db_pg_login"
	dbPGPasswordFlag = "db_pg_password"
	dbPGAddrFlag     = "db_pg_addr"
	dbPGNameFlag     = "db_pg_dbname"
	dbPGNoSSLFlag    = "db_pg_nossl"
	dbMigrationsFlag = "db_migrations"
)

var flags = []cli.Flag{
	cli.StringFlag{
		EnvVar: "CONFIG",
		Name:   configFlag,
		Value:  "config.yaml",
		Usage:  "Config file",
	},
	altsrc.NewStringFlag(cli.StringFlag{
		EnvVar: "PORT",
		Name:   portFlag,
		Value:  "8080",
		Usage:  "Application port",
	}),
	altsrc.NewBoolFlag(cli.BoolFlag{
		EnvVar: "DEBUG",
		Name:   debugFlag,
		Usage:  "start the server in debug mode",
	}),
	altsrc.NewBoolFlag(cli.BoolFlag{
		EnvVar: "TEXTLOG",
		Name:   textlogFlag,
		Usage:  "output log in text format",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		EnvVar: "PG_LOGIN",
		Name:   dbPGLoginFlag,
		Usage:  "DB Login (PostgreSQL)",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		EnvVar: "PG_PASSWORD",
		Name:   dbPGPasswordFlag,
		Usage:  "DB Password (PostgreSQL)",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		EnvVar: "PG_ADDR",
		Name:   dbPGAddrFlag,
		Usage:  "DB Address",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		EnvVar: "PG_DBNAME",
		Name:   dbPGNameFlag,
		Usage:  "DB name (PostgreSQL)",
	}),
	altsrc.NewBoolFlag(cli.BoolFlag{
		EnvVar: "PG_NOSSL",
		Name:   dbPGNoSSLFlag,
		Usage:  "DB disable SSL (PostgreSQL)",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		EnvVar: "MIGRATIONS_PATH",
		Name:   dbMigrationsFlag,
		Value:  "migrations",
		Usage:  "Location of DB migrations",
	}),
}

func setupLogs(c *cli.Context) {
	if c.Bool("debug") {
		gin.SetMode(gin.DebugMode)
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
		logrus.SetLevel(logrus.InfoLevel)
	}

	if c.Bool("textlog") {
		logrus.SetFormatter(&logrus.TextFormatter{})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}

func getDB(c *cli.Context) (db.DB, error) {
	url := fmt.Sprintf("postgres://%v:%v@%v/%v", c.String(dbPGLoginFlag), c.String(dbPGPasswordFlag), c.String(dbPGAddrFlag), c.String(dbPGNameFlag))
	if c.Bool(dbPGNoSSLFlag) {
		url = url + "?sslmode=disable"
	}
	return postgres.DBConnect(url, c.String(dbMigrationsFlag))
}
