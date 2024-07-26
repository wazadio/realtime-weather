package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/wazadio/realtime-weather/pkg/logger"
)

type DB struct {
	Read  *sql.DB
	Write *sql.DB
}

func NewDb(ctx context.Context) *DB {
	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	masterPort, err := strconv.Atoi(os.Getenv("DB_MASTER_PORT"))
	if err != nil {
		log.Fatalf("error connecting to master db : %s", err.Error())
	}

	masterConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_MASTER_HOST"), masterPort, os.Getenv("DB_MASTER_USERNAME"), os.Getenv("DB_MASTER_PASSWORD"), os.Getenv("DB_MASTER_NAME"))
	masterDb, err := sql.Open("postgres", masterConn)
	if err != nil {
		log.Fatalf("error connecting to master db : %s", err.Error())
	}

	err = masterDb.PingContext(dbCtx)
	if err != nil {
		log.Fatalf("error connecting to master db : %s", err.Error())
	}

	slavePort, err := strconv.Atoi(os.Getenv("DB_SLAVE_PORT"))
	if err != nil {
		log.Fatalf("error connecting to master db : %s", err.Error())
	}

	slaveConn := masterConn
	if os.Getenv("MASTER_DB_USERNAME") != "" && os.Getenv("MASTER_DB_PASSWORD") != "" {
		slaveConn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_SLAVE_HOST"), slavePort, os.Getenv("DB_SLAVE_USERNAME"), os.Getenv("DB_SLAVE_PASSWORD"), os.Getenv("DB_MASTER_NAME"))
	}

	slaveDb, err := sql.Open("postgres", slaveConn)
	if err != nil {
		log.Fatalf("error connecting to slave db : %s", err.Error())
	}

	err = slaveDb.PingContext(dbCtx)
	if err != nil {
		log.Fatalf("error connecting to salve db : %s", err.Error())
	}

	logger.Print(ctx, logger.INFO, "Database connected")

	return &DB{
		Read:  slaveDb,
		Write: masterDb,
	}
}
