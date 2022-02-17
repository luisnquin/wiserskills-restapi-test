package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/luisnquin/restapi-technical-test/src/controllers/tickets"
)

func ApplyTickets(g *echo.Group) {
	g.GET("s", tickets.FetchTickets())
	g.GET("/:id", tickets.FetchById())
	g.POST("", tickets.NewTicket())
	g.PATCH("/:id", tickets.ModifyTicketById())
	g.PUT("/:id", tickets.UpdateTicketById())
	g.DELETE("/:id", tickets.RemoveTicketById())
}
