package main

import (
	"viewExchange/app/controllers"
)

func main() {

	controllers.StreamIngestionData()

	controllers.StartWebServer()
}
