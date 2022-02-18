package tickets

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/luisnquin/restapi-technical-test/src/constants"
	"github.com/luisnquin/restapi-technical-test/src/models"
	"github.com/luisnquin/restapi-technical-test/src/storage"
)

func FetchTickets() echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			db  = storage.Get(constants.Persistence)
			err error
		)

		desc, _ := strconv.ParseBool(c.QueryParam("desc"))

		if err = db.Connect(); err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "tickets.get",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    500,
					Message: "Internal Server Error",
					Errors: []map[string]interface{}{
						{
							"reason":  "Internal Server Error",
							"message": "Database connection failed",
						},
					},
				},
			})
		}
		defer func() {
			if err = db.Close(); err != nil {
				panic(err)
			}
		}()

		var q string
		if desc {
			q = "SELECT * FROM tickets_view ORDER BY id DESC;"
		} else {
			q = "SELECT * FROM tickets_view;"
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		stmt, err := db.PrepareContext(ctx, q)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "tickets.get",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    500,
					Message: "Internal Server Error",
					Errors: []map[string]interface{}{
						{
							"reason":  "Internal Server Error",
						},
					},
				},
			})
		}
		defer func() {
			if err = stmt.Close(); err != nil {
				panic(err)
			}
		}()

		rows, err := stmt.QueryContext(ctx)
		if err != nil {
			return c.JSON(http.StatusNotFound, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "tickets.get",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    404,
					Message: "Not Found",
					Errors: []map[string]interface{}{
						{
							"reason":  "Not Found",
							"message": "There was an error when tried to bring the payload",
						},
					},
				},
			})
		}

		var tviews models.TicketViews

		for rows.Next() {
			var tview models.TicketView
			if err = rows.Scan(&tview.Id, &tview.Participant, &tview.Event); err != nil {
				return c.JSON(http.StatusConflict, models.BadResponse{
					APIVersion: constants.APIVersion,
					Method:     "tickets.get",
					Context:    c.Request().URL.String(),
					Error: models.Error{
						Code:    409,
						Message: "Conflict",
						Errors: []map[string]interface{}{
							{
								"reason": "Conflict",
								"message": "An error was logged while trying to process the payload",
							},
						},
					},
				})
			}
			tviews = append(tviews, tview)
		}
		if len(tviews) == 0 {
			return c.JSON(http.StatusNoContent, models.SuccessfulResponse{
				APIVersion: constants.APIVersion,
				Method:     "tickets.get",
				Context:    c.Request().URL.String(),
			})
		}
		return c.JSON(http.StatusOK, models.SuccessfulResponse{
			APIVersion: constants.APIVersion,
			Method:     "tickets.get",
			Context:    c.Request().URL.String(),
			Data:       tviews,
		})
	}
}

func FetchById() echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			db  = storage.Get(constants.Persistence)
			err error
		)

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "tickets.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    422,
					Message: "Unprocessable Entity",
					Errors: []map[string]interface{}{
						{
							"reason":  "Unprocessable Entity",
							"message": "The ID parameter cannot be processed as integer",
						},
					},
				},
			})
		}
		if err = db.Connect(); err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "tickets.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    500,
					Message: "Internal Server Error",
					Errors: []map[string]interface{}{
						{
							"reason":  "Internal Server Error",
							"message": "Database connection failed",
						},
					},
				},
			})
		}

		defer func() {
			if err = db.Close(); err != nil {
				panic(err)
			}
		}()

		var q string
		switch constants.Persistence {
		case storage.PostgreSQL:
			q = "SELECT * FROM tickets_view WHERE id = $1;"
		case storage.MySQL:
			q = "SELECT * FROM tickets_view WHERE id = ?;"
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		stmt, err := db.PrepareContext(ctx, q)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "tickets.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    500,
					Message: "Internal Server Error",
					Errors: []map[string]interface{}{
						{
							"reason":  "Internal Server Error",
						},
					},
				},
			})
		}
		defer func() {
			if err = stmt.Close(); err != nil {
				panic(err)
			}
		}()

		var tview models.TicketView
		err = stmt.QueryRowContext(ctx, id).Scan(&tview.Id, &tview.Participant, &tview.Event)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "tickets.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    400,
					Message: "Bad Request",
					Errors: []map[string]interface{}{
						{
							"reason":  "Bad Request",
							"message": "The ID paremeter was rejected, not valid",
						},
					},
				},
			})
		}

		return c.JSON(http.StatusOK, models.SuccessfulResponse{
			APIVersion: constants.APIVersion,
			Method:     "tickets.get",
			Context:    c.Request().URL.String(),
			Params: map[string]interface{}{
				"id": id,
			},
			Data: tview,
		})
	}
}
