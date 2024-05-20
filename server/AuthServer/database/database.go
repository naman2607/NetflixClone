// Database initialization and configuration code.

package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
	"github.com/naman2607/netflixClone/config"
)

var (
	db   *sql.DB
	mu   sync.Mutex
	once sync.Once
)

func GetDB() *sql.DB {
	mu.Lock()
	defer mu.Unlock()
	return db
}

func InitPostgresDB(dbConf *config.ServerConfig) error {

	once.Do(func() {
		postgresDb := dbConf.GetPostgresDBConf()
		psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", postgresDb.HOST, postgresDb.PORT, postgresDb.USER, postgresDb.PASSWORD, postgresDb.DBNAME)
		// open database
		var err error
		db, err = sql.Open("postgres", psqlconn)
		if err != nil {
			log.Fatal("Failed to open connection with postgres db ", err)
		}
		db.SetConnMaxIdleTime(10)
		db.SetMaxOpenConns(10)
		// check db
		err = db.Ping()
		if err != nil {
			log.Fatal("Failed to ping postgres db ", err)
		}
		log.Println("Postgres db Connected!")
	})
	return nil

}

func ExecuteTransactional(ctx context.Context, txFunc func(context.Context, *sql.Tx) error) (err error) {
	conn, err := db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection from pool: %w", err)
	}
	defer conn.Close()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = txFunc(ctx, tx)
	return err
}
