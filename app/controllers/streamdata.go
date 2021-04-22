package controllers

import (
	"log"
	"viewExchange/app/models"
	"viewExchange/config"
	"viewExchange/exchange/bitflyer"
	"viewExchange/exchange/gmocoin"
	"viewExchange/exchange/quoine"
)

func StreamIngestionData() {
	bitflyerChannel := make(chan bitflyer.Ticker)
	gmocoinChannel := make(chan gmocoin.Ticker)
	quoineChannel := make(chan quoine.Ticker)

	switch OpenExchange {
	case "bitflyer": //Bitflyer
		switch CurrencyPair {
		case "BTC_JPY":
			log.Println("BTC_JPY")
			go bitflyer.GetRealTimeTicker("BTC_JPY", bitflyerChannel)
		case "XRP_JPY":
			log.Println("XRP_JPY")
			go bitflyer.GetRealTimeTicker("XRP_JPY", bitflyerChannel)
		case "ETH_JPY":
			log.Println("ETH_JPY")
			go bitflyer.GetRealTimeTicker("ETH_JPY", bitflyerChannel)
		default:
			log.Println("BTC_JPY")
			go bitflyer.GetRealTimeTicker("BTC_JPY", bitflyerChannel)
		}
	case "gmocoin": //GMO Coin
		switch CurrencyPair {
		case "BTC_JPY":
			log.Println("BTC_JPY")
			go gmocoin.GetRealTimeTicker("BTC_JPY", gmocoinChannel)
		case "XRP_JPY":
			log.Println("XRP_JPY")
			go gmocoin.GetRealTimeTicker("XRP_JPY", gmocoinChannel)
		case "ETH_JPY":
			log.Println("ETH_JPY")
			go gmocoin.GetRealTimeTicker("ETH_JPY", gmocoinChannel)
		default:
			log.Println("BTC_JPY")
			go gmocoin.GetRealTimeTicker("BTC_JPY", gmocoinChannel)
		}
	case "quoine": //Liquid by Quoine
		switch CurrencyPair {
		case "BTC_JPY":
			log.Println("BTC_JPY")
			go quoine.GetRealTimeTicker("BTC_JPY", quoineChannel)
		case "XRP_JPY":
			log.Println("XRP_JPY")
			go quoine.GetRealTimeTicker("XRP_JPY", quoineChannel)
		case "ETH_JPY":
			log.Println("ETH_JPY")
			go quoine.GetRealTimeTicker("ETH_JPY", quoineChannel)
		default:
			log.Println("BTC_JPY")
			go quoine.GetRealTimeTicker("BTC_JPY", quoineChannel)
		}
	default: //Default Bitflyer
		switch CurrencyPair {
		case "BTC_JPY":
			log.Println("BTC_JPY")
			go bitflyer.GetRealTimeTicker("BTC_JPY", bitflyerChannel)
		case "XRP_JPY":
			log.Println("XRP_JPY")
			go bitflyer.GetRealTimeTicker("XRP_JPY", bitflyerChannel)
		case "ETH_JPY":
			log.Println("ETH_JPY")
			go bitflyer.GetRealTimeTicker("ETH_JPY", bitflyerChannel)
		default:
			log.Println("BTC_JPY")
			go bitflyer.GetRealTimeTicker("BTC_JPY", bitflyerChannel)
		}
	}

	go func() {
		switch OpenExchange {
		case "bitflyer":
			for ticker := range bitflyerChannel {
				log.Printf("Exchange: Bitflyer, Pair: %s, Timestamp: %s", ticker.ProductCode, ticker.Timestamp)
				//log.Printf("action=StreamIngestionData, %v", ticker)
				for _, duration := range config.Config.Durations {
					isCreated := models.BitflyerCreateCandles(ticker, ticker.ProductCode, duration)
					if isCreated == true {
						//log.Println("bitflyer_candle_created")
						//TODO
					}
				}
			}
		case "gmocoin":
			for ticker := range gmocoinChannel {
				log.Printf("Exchange: GMO Coin, Pair: %s, Timestamp: %s", ticker.Pair, ticker.Timestamp)
				//log.Printf("action=StreamIngestionData, %v", ticker)
				for _, duration := range config.Config.Durations {

					isCreated := models.GMOCoinCreateCandles(ticker, ticker.Pair, duration)
					if isCreated == true {
						//log.Println("gmocoin_candle_created")
						//TODO
					}
				}
			}
		case "quoine":
			for ticker := range quoineChannel {
				log.Printf("Exchange: Quoine, Pair: %s, Timestamp: %s", ticker.Pair, ticker.Timestamp)
				//log.Printf("action=StreamIngestionData, %v", ticker)
				for _, duration := range config.Config.Durations {
					isCreated := models.QuoineCreateCandles(ticker, ticker.Pair, duration)
					if isCreated == true {
						//log.Println("quoine_candle_created")
						//TODO
					}
				}
			}
		default: //Default Bitflyer
			for ticker := range bitflyerChannel {
				log.Printf("Exchange: Bitflyer, Pair: %s, Timestamp: %s", ticker.ProductCode, ticker.Timestamp)
				//log.Printf("action=StreamIngestionData, %v", ticker)
				for _, duration := range config.Config.Durations {
					isCreated := models.BitflyerCreateCandles(ticker, ticker.ProductCode, duration)
					if isCreated == true {
						//log.Println("bitflyer_candle_created")
						//TODO
					}
				}
			}
		}
	}()
}
