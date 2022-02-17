package routers

import "github.com/labstack/echo/v4"

func Apply(e *echo.Echo) {
	persistence := e.Group("/persistence")
	ApplyPersistence(persistence)

	api := e.Group("/api")
	v1 := api.Group("/v1")

	event := v1.Group("/event")
	ApplyEvents(event)

	participant := v1.Group("/participant")
	ApplyParticipants(participant)

	ticket := v1.Group("/ticket")
	ApplyTickets(ticket)
}
