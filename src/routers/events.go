package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/luisnquin/restapi-technical-test/src/controllers/events"
)

func ApplyEvents(g *echo.Group) {
	g.GET("s", events.Fetch())
	g.GET("/:id", events.ById())
	g.GET("/:id/tickets", events.FetchTicketsById())
	g.GET("/:id/participants", events.FetchParticipantsById())
	g.GET("/:event-id/participant/:participant-id", events.FetchParticipantByIds())
	g.POST("", events.New())
	g.POST("/:event-id/participant/:participant-id", events.NewParticipantByIds())
	g.PUT("/:id", events.UpdateById())
	g.DELETE("/:id", events.RemoveById())
	g.DELETE("/:id/participants", events.RemoveByIdWithParticipants())
}
