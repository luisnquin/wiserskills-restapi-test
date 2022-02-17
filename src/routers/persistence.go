package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/luisnquin/restapi-technical-test/src/controllers/persistence"
)

func ApplyPersistence(g *echo.Group) {
	g.GET("/help", persistence.Help())
	g.POST("/build/:persistence", persistence.Build())
}
