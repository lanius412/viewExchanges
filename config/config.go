package config

import (
	"log"
	"time"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	BitflyerDbName string
	GMOCoinDbName  string
	QuoineDbName   string
	SQLDriver      string

	Durations []time.Duration
	Pairs     []string
}

var Config ConfigList

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("Failed to Read Config File: %v", err)
	}

	durations := []time.Duration{
		time.Minute * 1,
		time.Minute * 5,
		time.Minute * 30,
		time.Hour * 1,
	}

	pairs := []string{
		"BTC_JPY",
		"ETH_JPY",
		"XRP_JPY",
	}

	Config = ConfigList{
		BitflyerDbName: cfg.Section("db").Key("bitflyerDbName").String(),
		GMOCoinDbName:  cfg.Section("db").Key("gmocoinDbName").String(),
		QuoineDbName:   cfg.Section("db").Key("quoineDbName").String(),

		SQLDriver: cfg.Section("db").Key("driver").String(),
		Durations: durations,
		Pairs:     pairs,
	}
}
