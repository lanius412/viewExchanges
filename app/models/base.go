package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"viewExchange/config"
)

var TableNames = make([]string, 0)
var BitflyerDbConnection, GMOCoinDbConnection, QuoineDbConnection *sql.DB

func GetCandleTableName(pair string, duration time.Duration) string {
	return fmt.Sprintf("%s_%s", pair, duration)
}

func init() {
	var err error
	var cfg = config.Config

	for _, pair := range cfg.Pairs {
		for _, duration := range cfg.Durations {
			TableNames = append(TableNames, GetCandleTableName(pair, duration))
		}
	}

	BitflyerDbConnection, err = sql.Open(cfg.SQLDriver, cfg.BitflyerDbName)
	if err != nil {
		log.Fatalln(err)
	}
	for _, table := range TableNames {
		c := fmt.Sprintf(`
					CREATE TABLE IF NOT EXISTS %s (
						time DATETIME PRIMARY KEY NOT NULL,
						open FLOAT,
						close FLOAT,
						high FLOAT,
						low FLOAT,
						volume FLOAT
					)`, table)
		BitflyerDbConnection.Exec(c)
	}

	GMOCoinDbConnection, err = sql.Open(cfg.SQLDriver, cfg.GMOCoinDbName)
	if err != nil {
		log.Fatalln(err)
	}
	for _, table := range TableNames {
		c := fmt.Sprintf(`
						CREATE TABLE IF NOT EXISTS %s (
							time DATETIME PRIMARY KEY NOT NULL,
							open FLOAT,
							close FLOAT,
							high FLOAT,
							low FLOAT,
							volume FLOAT
						)`, table)
		GMOCoinDbConnection.Exec(c)
	}

	QuoineDbConnection, err = sql.Open(cfg.SQLDriver, cfg.QuoineDbName)
	if err != nil {
		log.Fatalln(err)
	}
	for _, table := range TableNames {
		c := fmt.Sprintf(`
					CREATE TABLE IF NOT EXISTS %s (
						time DATETIME PRIMARY KEY NOT NULL,
						open FLOAT,
						close FLOAT,
						high FLOAT,
						low FLOAT,
						volume FLOAT
					)`, table)
		QuoineDbConnection.Exec(c)
	}

}
