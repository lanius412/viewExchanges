package models

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
	"viewExchange/exchange/bitflyer"
	"viewExchange/exchange/gmocoin"
	"viewExchange/exchange/quoine"
)

type Candle struct {
	Pair     string        `json:"pair"`
	Duration time.Duration `json:"duration"`
	Time     time.Time     `json:"time"`
	Open     float64       `json:"open"`
	Close    float64       `json:"close"`
	High     float64       `json:"high"`
	Low      float64       `json:"low"`
	Volume   float64       `json:"volume"`
}

func NewCandle(pair string, duration time.Duration, dateTime time.Time, open, close, high, low, volume float64) *Candle {
	return &Candle{
		pair,
		duration,
		dateTime,
		open,
		close,
		high,
		low,
		volume,
	}
}

func (c *Candle) TableName() string {
	return GetCandleTableName(c.Pair, c.Duration)
}

func (c *Candle) Create(exchange string) (err error) {
	cmd := fmt.Sprintf(`
		INSERT INTO %s (
			time,
			open,
			close,
			high,
			low,
			volume) VALUES (?, ?, ?, ?, ?, ?)`, c.TableName())

	switch exchange {
	case "bitflyer":
		_, err = BitflyerDbConnection.Exec(cmd, c.Time.Format(time.RFC3339), c.Open, c.Close, c.High, c.Low, c.Volume)
		if err != nil {
			return err
		}
	case "gmocoin":
		_, err = GMOCoinDbConnection.Exec(cmd, c.Time.Format(time.RFC3339), c.Open, c.Close, c.High, c.Low, c.Volume)
		if err != nil {
			return err
		}
	case "quoine":
		_, err = QuoineDbConnection.Exec(cmd, c.Time.Format(time.RFC3339), c.Open, c.Close, c.High, c.Low, c.Volume)
		if err != nil {
			return err
		}
	}

	return err
}

func (c *Candle) Save(exchange string) (err error) {
	cmd := fmt.Sprintf(`
		UPDATE %s SET 
			open = ?,
			close = ?,
			high = ?,
			low = ?,
			volume = ? WHERE time = ?`, c.TableName())

	switch exchange {
	case "bitflyer":
		_, err = BitflyerDbConnection.Exec(cmd, c.Open, c.Close, c.High, c.Low, c.Volume, c.Time.Format(time.RFC3339))
		if err != nil {
			return err
		}
	case "gmocoin":
		_, err = GMOCoinDbConnection.Exec(cmd, c.Open, c.Close, c.High, c.Low, c.Volume, c.Time.Format(time.RFC3339))
		if err != nil {
			return err
		}
	case "quoine":
		_, err = QuoineDbConnection.Exec(cmd, c.Open, c.Close, c.High, c.Low, c.Volume, c.Time.Format(time.RFC3339))
		if err != nil {
			return err
		}
	}

	return err
}

func GetCandle(exchange, pair string, duration time.Duration, dateTime time.Time) *Candle {
	tableName := GetCandleTableName(pair, duration)
	cmd := fmt.Sprintf(`
		SELECT time, open, close, high, low, volume FROM %s WHERE time = ?`, tableName)

	var row *sql.Row
	switch exchange {
	case "bitflyer":
		row = BitflyerDbConnection.QueryRow(cmd, dateTime.Format(time.RFC3339))
	case "gmocoin":
		row = GMOCoinDbConnection.QueryRow(cmd, dateTime.Format(time.RFC3339))
	case "quoine":
		row = QuoineDbConnection.QueryRow(cmd, dateTime.Format(time.RFC3339))
	}
	var candle Candle
	err := row.Scan(
		&candle.Time,
		&candle.Open,
		&candle.Close,
		&candle.High,
		&candle.Low,
		&candle.Volume)
	if err != nil {
		return nil
	}
	return NewCandle(pair, duration, candle.Time, candle.Open, candle.Close, candle.High, candle.Low, candle.Volume)
}

func BitflyerCreateCandles(ticker bitflyer.Ticker, pair string, duration time.Duration) bool {
	currentCandle := GetCandle("bitflyer", pair, duration, ticker.TruncateDateTime(duration))
	price := ticker.GetMidPrice()
	if currentCandle == nil {
		candle := NewCandle(pair, duration, ticker.TruncateDateTime(duration), price, price, price, price, ticker.Volume)
		candle.Create("bitflyer")
		return true
	}
	if currentCandle.High <= price {
		currentCandle.High = price
	} else if currentCandle.Low >= price {
		currentCandle.Low = price
	}
	currentCandle.Volume += ticker.Volume
	currentCandle.Close = price
	currentCandle.Save("bitflyer")

	return false
}

func GMOCoinCreateCandles(ticker gmocoin.Ticker, pair string, duration time.Duration) bool {
	currentCandle := GetCandle("gmocoin", pair, duration, ticker.TruncateDateTime(duration))
	price := ticker.GetMidPrice()
	volume, err := strconv.ParseFloat(ticker.Volume, 64)
	if err != nil {
		log.Fatalln(err)
	}
	if currentCandle == nil {
		candle := NewCandle(pair, duration, ticker.TruncateDateTime(duration), price, price, price, price, volume)
		candle.Create("gmocoin")
		return true
	}
	if currentCandle.High <= price {
		currentCandle.High = price
	} else if currentCandle.Low >= price {
		currentCandle.Low = price
	}
	currentCandle.Volume += volume
	currentCandle.Close = price
	currentCandle.Save("gmocoin")

	return false
}

func QuoineCreateCandles(ticker quoine.Ticker, pair string, duration time.Duration) bool {
	currentCandle := GetCandle("quoine", pair, duration, ticker.TruncateDateTime(duration))
	price := ticker.GetMidPrice()
	volume, err := strconv.ParseFloat(ticker.Volume, 64)
	if err != nil {
		log.Fatalln(err)
	}
	if currentCandle == nil {
		candle := NewCandle(pair, duration, ticker.TruncateDateTime(duration), price, price, price, price, volume)
		candle.Create("quoine")
		return true
	}
	if currentCandle.High <= price {
		currentCandle.High = price
	} else if currentCandle.Low >= price {
		currentCandle.Low = price
	}
	currentCandle.Volume += volume
	currentCandle.Close = price
	currentCandle.Save("quoine")

	return false
}

func GetAllCandle(exchange, pair string, duration time.Duration, limit int) (dfCandle *DataFrameCandle, err error) {
	tableName := GetCandleTableName(pair, duration)
	cmd := fmt.Sprintf(`SELECT * FROM (
		SELECT time, open, close, high, low, volume FROM %s ORDER BY time DESC LIMIT ?
		) ORDER BY time ASC;`, tableName)

	var rows *sql.Rows
	switch exchange {
	case "bitflyer":
		rows, err = BitflyerDbConnection.Query(cmd, limit)
	case "gmocoin":
		rows, err = GMOCoinDbConnection.Query(cmd, limit)
	case "quoine":
		rows, err = QuoineDbConnection.Query(cmd, limit)
	}
	if err != nil {
		log.Fatalln("DB Query error:", err)
	}
	defer rows.Close()

	dfCandle = &DataFrameCandle{}
	dfCandle.Pair = pair
	dfCandle.Duration = duration
	for rows.Next() {
		var candle Candle
		candle.Pair = pair
		candle.Duration = duration
		rows.Scan(
			&candle.Time,
			&candle.Open,
			&candle.Close,
			&candle.High,
			&candle.Low,
			&candle.Volume)
		dfCandle.Candles = append(dfCandle.Candles, candle)
	}
	err = rows.Err()
	if err != nil {
		return
	}
	return dfCandle, nil
}
