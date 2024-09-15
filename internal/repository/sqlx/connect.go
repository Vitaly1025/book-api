package sqlx

import (
	"context"
	"fmt"
	"log"
	"time"

	"book-api/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // <------------ here
)

func NewPostgresDB(cfg config.DataBase) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Name, cfg.Password, cfg.SSL)

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	go func() {
		for {
			err := db.PingContext(context.Background())
			if err != nil {
				log.Println("Lost connection to the database. Reconnecting...")
				db, _ = sqlx.Connect("postgres", connectionString)
			}
			time.Sleep(10 * time.Second)
		}
	}()

	if err != nil {
		return nil, err
	}

	return db, nil
}
