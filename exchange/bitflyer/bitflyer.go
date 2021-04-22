package bitflyer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type Ticker struct {
	ProductCode string  `json:"product_code"`
	Timestamp   string  `json:"timestamp"`
	BestBid     float64 `json:"best_bid"`
	BestAsk     float64 `json:"best_ask"`
	Volume      float64 `json:"volume"`
}

func (t *Ticker) GetMidPrice() float64 {
	return (t.BestBid + t.BestAsk) / 2
}

func (t *Ticker) DateTime() time.Time {
	dateTime, err := time.Parse(time.RFC3339, t.Timestamp)
	if err != nil {
		log.Printf("action=ParseDateTime err=%s", err.Error())
	}
	return dateTime
}

func (t *Ticker) TruncateDateTime(duration time.Duration) time.Time {
	return t.DateTime().Truncate(duration)
}

type JsonRPC2 struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Result  interface{} `json:"result,omitempty"`
	ID      *int        `json:"id,omitempty"`
}
type SubscribeParams struct {
	Channel string `json:"channel"`
}

func GetRealTimeTicker(pair string, ch chan<- Ticker) {
	u := url.URL{Scheme: "wss", Host: "ws.lightstream.bitflyer.com", Path: "/json-rpc"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalln("dial:", err)
	}
	defer c.Close()

	channel := fmt.Sprintf("lightning_ticker_%s", pair)
	if err := c.WriteJSON(&JsonRPC2{Version: "2.0", Method: "subscribe", Params: &SubscribeParams{channel}}); err != nil {
		log.Fatalln("write:", err)
		return
	}
	c.SetWriteDeadline(time.Now().Add(10 * time.Second))

	for {
		message := new(JsonRPC2)
		if err := c.ReadJSON(message); err != nil {
			log.Fatalln("read:", err)
			return
		}
		c.SetReadDeadline(time.Now().Add(10 * time.Second))
		if message.Method == "channelMessage" {
			switch v := message.Params.(type) {
			case map[string]interface{}:
				for key, value := range v {
					if key == "message" {
						marshalTic, err := json.Marshal(value)
						if err != nil {
							log.Fatalln(err)
						}
						var ticker Ticker
						if err := json.Unmarshal(marshalTic, &ticker); err != nil {
							log.Fatalln(err)
						}
						ch <- ticker
					}
				}
			}
		}
	}
	return
}
