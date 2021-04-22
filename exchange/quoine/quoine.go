package quoine

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Ticker struct {
	Pair      string `json:"currency_pair_code"`
	Ask       string `json:"market_ask"`
	Bid       string `json:"market_bid"`
	Timestamp string `json:"timestamp"`
	Volume    string `json:"volume_24h"`
}

func (t *Ticker) GetMidPrice() float64 {
	bid, err := strconv.ParseFloat(t.Bid, 64)
	ask, err := strconv.ParseFloat(t.Ask, 64)
	if err != nil {
		log.Printf("action=ParseFloat err=%s", err.Error())
	}
	return (bid + ask) / 2
}

func (t *Ticker) DateTime() time.Time {
	timestamp, err := strconv.ParseFloat(t.Timestamp, 64)
	if err != nil {
		log.Fatalln("parse timestamp error", err)
	}
	sec, nsec := math.Modf(timestamp)
	dateTime := time.Unix(int64(sec), int64(nsec*(1e9)))
	return dateTime
}

func (t *Ticker) TruncateDateTime(duration time.Duration) time.Time {
	return t.DateTime().Truncate(duration)
}

type sendMsg struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}
type subscriveData struct {
	Channel string `json:"channel"`
}

var pair_id = map[string]string{
	"btcjpy":  "5",
	"ethjpy":  "29",
	"xrpjpy":  "83",
	"bchjpy":  "41",
	"qashjpy": "50",
}

func GetRealTimeTicker(pair string, ch chan<- Ticker) {
	u := url.URL{Scheme: "wss", Host: "tap.liquid.com", Path: "/app/LiquidTapClient"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalln("dial:", err)
	}
	defer c.Close()

	pairLower := strings.ToLower(strings.Replace(pair, "_", "", 1))
	channel := fmt.Sprintf("product_cash_%s_%s", pairLower, pair_id[pairLower])
	if err := c.WriteJSON(&sendMsg{Event: "pusher:subscribe", Data: &subscriveData{Channel: channel}}); err != nil {
		log.Fatalln("write:", err)
	}
	c.SetWriteDeadline(time.Now().Add(10 * time.Second))

	for {
		var message interface{}
		if err := c.ReadJSON(&message); err != nil {
			log.Fatalln("read:", err)
		}
		c.SetReadDeadline(time.Now().Add(10 * time.Second))
		switch v := message.(type) {
		case map[string]interface{}:
			if v["event"] == "updated" {
				data := v["data"].(string)
				var ticker Ticker
				if err := json.Unmarshal([]byte(data), &ticker); err != nil {
					log.Fatalln("unmarshal:", err)
				}
				ch <- ticker
			} else {
				return
			}
		}
	}
}
