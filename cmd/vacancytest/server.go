package main

import (
	"fmt"
	"os"

	"time"

	"context"
	"net/http"
	"os/signal"

	"text/tabwriter"

	"github.com/Iliad/vacancytest/pkg/router"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
)

func initServer(c *cli.Context) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent|tabwriter.Debug)
	for _, f := range c.GlobalFlagNames() {
		fmt.Fprintf(w, "Flag: %s\t Value: %s\n", f, c.String(f))
	}
	w.Flush()

	setupLogs(c)

	db, err := getDB(c)
	exitOnErr(err)
	defer db.Close()

	app := router.CreateRouter(&db)

	// graceful shutdown support
	srv := http.Server{
		Addr:    ":" + c.String(portFlag),
		Handler: app,
	}

	go exitOnErr(srv.ListenAndServe())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logrus.Infoln("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
