package main

import (
	"github.com/catfishlty/webhooks-hub/internal/api"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	g errgroup.Group
)

func main() {
	sqlitePath := "test.db"
	port := 9102
	db, err := gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	hub := api.NewHub(db)
	hub.Migrate()
	g.Go(func() error {
		return hub.Server(port).ListenAndServe()
	})
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
