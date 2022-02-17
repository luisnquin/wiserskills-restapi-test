package main

import (
	"fmt"
	"time"

	"github.com/TwiN/go-color"
	"github.com/labstack/echo/v4"

	"github.com/luisnquin/restapi-technical-test/src/middleware"
	"github.com/luisnquin/restapi-technical-test/src/routers"
)

func main() {
	var server = echo.New()
	middleware.Apply(server)
	routers.Apply(server)

	go func() {
		time.Sleep(time.Millisecond * 250)
		fmt.Printf("Fast acccess:\n %s\n %s\n %s\n\n",
			color.InCyan(" -> http://127.0.0.1:8000/api/v1/participants"),
			color.InPurple(" -> http://127.0.0.1:8000/api/v1/tickets"),
			color.InYellow(" -> http://127.0.0.1:8000/api/v1/events"),
		)
		fmt.Printf("First access:\n %s\n\n",
			color.InRed(" -> http://127.0.0.1:8000/persistence/help"),
		)
	}()
	server.Logger.Fatal(server.Start(":8000"))
}
