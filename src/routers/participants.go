package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/luisnquin/restapi-technical-test/src/controllers/participants"
)

func ApplyParticipants(g *echo.Group) {
	g.GET("s", participants.Fetch())
	g.GET("/:id", participants.FetchById())
	g.GET("/:id/tickets", participants.FetchTicketsById())
	g.POST("", participants.New())
	g.PUT("/:id", participants.UpdateById())
	g.DELETE("/:id", participants.RemoveById())
}
