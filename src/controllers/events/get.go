package events

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"github.com/luisnquin/restapi-technical-test/src/constants"
	"github.com/luisnquin/restapi-technical-test/src/models"
	"github.com/luisnquin/restapi-technical-test/src/storage"
)

func Fetch() echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			db  = storage.Get(constants.Persistence)
			err error
		)

		if err = db.Connect(); err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    500,
					Message: "Internal server error",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Internal server error",
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

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		stmt, err := db.PrepareContext(ctx, "SELECT * FROM events;")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    500,
					Message: "Internal server error",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Internal server error",
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
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    500,
					Message: "Internal server error",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Internal server error",
						},
					},
				},
			})
		}

		defer func() {
			if err = rows.Close(); err != nil {
				panic(err)
			}
		}()

		var events models.Events
		for rows.Next() {
			var e models.Event

			if err = rows.Scan(&e.Id, &e.Name, &e.Created_at); err != nil {
				return c.JSON(http.StatusConflict, models.BadResponse{
					APIVersion: constants.APIVersion,
					Method:     "events.get",
					Context:    c.Request().URL.String(),
					Error: models.Error{
						Code:    409,
						Message: "Conflict",
						Errors: []map[string]interface{}{
							{
								"reason":  err,
								"message": "Conflict",
							},
						},
					},
				})
			}
			events = append(events, e)
		}

		if len(events) == 0 {
			return c.JSON(http.StatusNoContent, models.SuccessfulResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
			})
		}

		return c.JSON(http.StatusOK, models.SuccessfulResponse{
			APIVersion: constants.APIVersion,
			Method:     "events.get",
			Context:    c.Request().URL.String(),
			Data:       events,
		})
	}
}

func ById() echo.HandlerFunc {
	return func(c echo.Context) error {
		var db = storage.Get(constants.Persistence)

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    422,
					Message: "Unprocessable entity",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Unprocessable entity",
						},
					},
				},
			})
		}
		if err = db.Connect(); err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    500,
					Message: "Internal server error",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Internal server error",
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

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		q := "SELECT * FROM events WHERE id = ? LIMIT 1;"
		if constants.Persistence == storage.PostgreSQL {
			q = sqlx.Rebind(sqlx.DOLLAR, q)
		}

		// Ask me again
		stmt, err := db.PrepareContext(ctx, q)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    500,
					Message: "Internal server error",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Internal server error",
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

		var event models.Event

		err = stmt.QueryRowContext(ctx, id).Scan(&event.Id, &event.Name, &event.Created_at)
		if err != nil {
			return c.JSON(http.StatusNotFound, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    404,
					Message: "Not Found",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Not Found",
						},
					},
				},
			})
		}

		return c.JSON(http.StatusOK, models.SuccessfulResponse{
			APIVersion: constants.APIVersion,
			Method:     "events.get",
			Context:    c.Request().URL.String(),
			Params: map[string]interface{}{
				"id": id,
			},
			Data: event,
		})
	}
}

func FetchTicketsById() echo.HandlerFunc {
	return func(c echo.Context) error {
		var db = storage.Get(constants.Persistence)

		desc, err := strconv.ParseBool(c.QueryParam("desc"))

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    422,
					Message: "Unprocessable entity",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Unprocessable entity",
						},
					},
				},
			})
		}

		if err = db.Connect(); err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    500,
					Message: "Internal server error",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Internal server error",
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
		switch {
		case constants.Persistence == storage.PostgreSQL && desc:
			q = "SELECT t.id AS id, CONCAT(p.firstname, ' ',p.lastname) AS participant, e.name AS event FROM tickets AS t INNER JOIN events AS e ON e.id=t.event INNER JOIN participants AS p ON p.id=t.participant WHERE e.id = $1 ORDER BY t.id DESC;"
		case constants.Persistence == storage.MySQL && desc:
			q = "SELECT t.id AS id, CONCAT(p.firstname, ' ',p.lastname) AS participant, e.name AS event FROM tickets AS t INNER JOIN events AS e ON e.id=t.event INNER JOIN participants AS p ON p.id=t.participant WHERE e.id = ? ORDER BY t.id DESC;"
		case constants.Persistence == storage.PostgreSQL:
			q = "SELECT t.id AS id, CONCAT(p.firstname, ' ',p.lastname) AS participant, e.name AS event FROM tickets AS t INNER JOIN events AS e ON e.id=t.event INNER JOIN participants AS p ON p.id=t.participant WHERE e.id = $1 ORDER BY t.id ASC;"
		case constants.Persistence == storage.MySQL:
			q = "SELECT t.id AS id, CONCAT(p.firstname, ' ',p.lastname) AS participant, e.name AS event FROM tickets AS t INNER JOIN events AS e ON e.id=t.event INNER JOIN participants AS p ON p.id=t.participant WHERE e.id = ? ORDER BY t.id ASC;"
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		stmt, err := db.PrepareContext(ctx, q)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    500,
					Message: "Internal server error",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Internal server error",
						},
					},
				},
			})
		}
		defer func() {
			if err != nil {
				panic(err)
			}
		}()
		rows, err := stmt.QueryContext(ctx, id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    500,
					Message: "Internal server error",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Internal server error",
						},
					},
				},
			})
		}
		defer func() {
			if err = rows.Close(); err != nil {
				panic(err)
			}
		}()

		var tviews models.TicketViews

		for rows.Next() {
			var tview models.TicketView

			if err = rows.Scan(&tview.Id, &tview.Participant, &tview.Event); err != nil {
				return c.JSON(http.StatusConflict, models.BadResponse{
					APIVersion: constants.APIVersion,
					Method:     "events.get",
					Context:    c.Request().URL.String(),
					Params: map[string]interface{}{
						"id": id,
					},
					Error: models.Error{
						Code:    409,
						Message: "Conflict",
						Errors: []map[string]interface{}{
							{
								"reason":  err,
								"message": "Conflict",
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
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
			})
		}
		return c.JSON(http.StatusOK, models.SuccessfulResponse{
			APIVersion: constants.APIVersion,
			Method:     "events.get",
			Context:    c.Request().URL.String(),
			Params: map[string]interface{}{
				"id": id,
			},
			Data: tviews,
		})
	}
}

func FetchParticipantsById() echo.HandlerFunc {
	return func(c echo.Context) error {
		var db = storage.Get(constants.Persistence)

		desc, err := strconv.ParseBool(c.QueryParam("desc"))

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    422,
					Message: "Unprocessable entity",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Unprocessable entity",
						},
					},
				},
			})
		}

		if err = db.Connect(); err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    500,
					Message: "Internal server error",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Internal server error",
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
		switch {
		case constants.Persistence == storage.PostgreSQL && desc:
			q = "SELECT t.id AS id, CONCAT(p.firstname, ' ',p.lastname) AS participant, e.name AS event FROM tickets AS t INNER JOIN events AS e ON e.id=t.event INNER JOIN participants AS p ON p.id=t.participant WHERE e.id = $1 ORDER BY t.id DESC;"
		case constants.Persistence == storage.MySQL && desc:
			q = "SELECT t.id AS id, CONCAT(p.firstname, ' ',p.lastname) AS participant, e.name AS event FROM tickets AS t INNER JOIN events AS e ON e.id=t.event INNER JOIN participants AS p ON p.id=t.participant WHERE e.id = ? ORDER BY t.id DESC;"
		case constants.Persistence == storage.PostgreSQL:
			q = "SELECT t.id AS id, CONCAT(p.firstname, ' ',p.lastname) AS participant, e.name AS event FROM tickets AS t INNER JOIN events AS e ON e.id=t.event INNER JOIN participants AS p ON p.id=t.participant WHERE e.id = $1 ORDER BY t.id ASC;"
		case constants.Persistence == storage.MySQL:
			q = "SELECT t.id AS id, CONCAT(p.firstname, ' ',p.lastname) AS participant, e.name AS event FROM tickets AS t INNER JOIN events AS e ON e.id=t.event INNER JOIN participants AS p ON p.id=t.participant WHERE e.id = ? ORDER BY t.id ASC;"
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		stmt, err := db.PrepareContext(ctx, q)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    500,
					Message: "Internal server error",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Internal server error",
						},
					},
				},
			})
		}
		defer func() {
			if err != nil {
				panic(err)
			}
		}()
		rows, err := stmt.QueryContext(ctx, id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    500,
					Message: "Internal server error",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Internal server error",
						},
					},
				},
			})
		}
		defer func() {
			if err = rows.Close(); err != nil {
				panic(err)
			}
		}()

		var tviews models.TicketViews

		for rows.Next() {
			var tview models.TicketView

			if err = rows.Scan(&tview.Id, &tview.Participant, &tview.Event); err != nil {
				return c.JSON(http.StatusConflict, models.BadResponse{
					APIVersion: constants.APIVersion,
					Method:     "events.get",
					Context:    c.Request().URL.String(),
					Params: map[string]interface{}{
						"id": id,
					},
					Error: models.Error{
						Code:    409,
						Message: "Conflict",
						Errors: []map[string]interface{}{
							{
								"reason":  err,
								"message": "Conflict",
							},
						},
					},
				})
			}
			tviews = append(tviews, tview)
		}
		if len(tviews) == 0 {
			return c.JSON(http.StatusNotFound, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    404,
					Message: "Not found",
					Errors: []map[string]interface{}{
						{
							"reason":  "Either the event does not exist or it has no participants",
							"message": "Not found",
						},
					},
				},
			})
		}
		return c.JSON(http.StatusOK, models.SuccessfulResponse{
			APIVersion: constants.APIVersion,
			Method:     "events.get",
			Context:    c.Request().URL.String(),
			Params: map[string]interface{}{
				"id": id,
			},
			Data: tviews,
		})
	}
}

func FetchParticipantByIds() echo.HandlerFunc {
	return func(c echo.Context) error {
		var db = storage.Get(constants.Persistence)

		event_id, err := strconv.Atoi(c.Param("event-id"))
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"event_id":       event_id,
					"participant_id": "",
				},
				Error: models.Error{
					Code:    422,
					Message: "Unprocessable entity",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Unprocessable entity",
						},
					},
				},
			})
		}

		participant_id, err := strconv.Atoi(c.Param("participant-id"))
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"event_id":       event_id,
					"participant_id": participant_id,
				},
				Error: models.Error{
					Code:    422,
					Message: "Unprocessable entity",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Unprocessable entity",
						},
					},
				},
			})
		}

		if err = db.Connect(); err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"event_id":       event_id,
					"participant_id": participant_id,
				},
				Error: models.Error{
					Code:    500,
					Message: "Internal server error",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Internal server error",
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
			q = "SELECT t.id AS id, CONCAT(p.firstname, ' ',p.lastname) AS participant, e.name AS event FROM tickets AS t INNER JOIN events AS e ON e.id=t.event INNER JOIN participants AS p ON p.id=t.participant WHERE e.id = $1 AND p.id = $2;"

		case storage.MySQL:
			q = "SELECT t.id AS id, CONCAT(p.firstname, ' ',p.lastname) AS participant, e.name AS event FROM tickets AS t INNER JOIN events AS e ON e.id=t.event INNER JOIN participants AS p ON p.id=t.participant WHERE e.id = ? AND p.id = ?;"
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		stmt, err := db.PrepareContext(ctx, q)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"event_id":       event_id,
					"participant_id": participant_id,
				},
				Error: models.Error{
					Code:    500,
					Message: "Internal server error",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Internal server error",
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
		err = stmt.QueryRowContext(ctx, event_id, participant_id).Scan(&tview.Id, &tview.Participant, &tview.Event)
		if err != nil {
			return c.JSON(http.StatusNotFound, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"event_id":       event_id,
					"participant_id": participant_id,
				},
				Error: models.Error{
					Code:    404,
					Message: "Not Found",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Not Found",
						},
					},
				},
			})
		}

		return c.JSON(http.StatusOK, models.SuccessfulResponse{
			APIVersion: constants.APIVersion,
			Method:     "events.get",
			Context:    c.Request().URL.String(),
			Params: map[string]interface{}{
				"event_id":       event_id,
				"participant_id": participant_id,
			},
			Data: tview,
		})
	}
}
