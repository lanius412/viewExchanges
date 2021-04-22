package gmocoin

import (
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type Ticker struct {
	Ask       string    `json:"ask"`
	Bid       string    `json:"bid"`
	Pair      string    `json:"symbol"`
	Timestamp time.Time `json:"timestamp"`
	Volume    string    `json:"volume"`
}

func (t *Ticker) GetMidPrice() float64 {
	bid, err := strconv.ParseFloat(t.Bid, 64)
	ask, err := strconv.ParseFloat(t.Ask, 64)
	if err != nil {
		log.Printf("action=ParseFloat err=%s", err.Error())
	}
	return (bid + ask) / 2
}

func (t *Ticker) TruncateDateTime(duration time.Duration) time.Time {
	return t.Timestamp.Truncate(duration)
}

type sendMsg struct {
	Command string `json:"command"`
	Channel string `json:"channel"`
	Symbol  string `json:"symbol"`
}

func GetRealTimeTicker(pair string, ch chan<- Ticker) {

	u := url.URL{Scheme: "wss", Host: "api.coin.z.com", Path: "/ws/public/v1"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalln("dial:", err)
	}
	defer c.Close()

	if err := c.WriteJSON(&sendMsg{Command: "subscribe", Channel: "ticker", Symbol: pair}); err != nil {
		log.Fatalln("write:", err)
	}
	c.SetWriteDeadline(time.Now().Add(10 * time.Second))

	for {
		var message map[string]string
		if err := c.ReadJSON(&message); err != nil {
			log.Fatalln("read:", err)
		}
		c.SetReadDeadline(time.Now().Add(10 * time.Second))
		for key, _ := range message {
			if key == "error" {
				return
			} else {
				timeUTC, err := time.Parse(time.RFC3339, message["timestamp"])
				if err != nil {
					log.Println(err)
				}
				timeJST := timeUTC.Add(9 * time.Hour)
				ticker := Ticker{
					Ask:       message["ask"],
					Bid:       message["bid"],
					Pair:      message["symbol"],
					Timestamp: timeJST,
					Volume:    message["volume"],
				}
				ch <- ticker
				break
			}
		}
	}

}
