package events

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

func Fetch() echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			db  = storage.Get(constants.Persistence)
			err error
		)

		desc, _ := strconv.ParseBool(c.QueryParam("desc"))

		if err = db.Connect(); err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
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

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		var q string
		if desc {
			q = "SELECT * FROM events ORDER BY id DESC;"
		} else {
			q = "SELECT * FROM events;"
		}

		stmt, err := db.PrepareContext(ctx, q)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    500,
					Message: "Internal Server Error",
					Errors: []map[string]interface{}{
						{
							"reason": "Internal Server Error",
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
							"reason":  "Internal Server Error",
							"message": "There was an error when tried to bring the payload",
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
								"reason":  "Conflict",
								"message": "An error was logged while trying to process the payload",
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
							"reason":  "Unprocessable Entity",
							"message": "The provided ID parameter cannot be processed as integer",
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

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		var q string
		switch {
		case constants.Persistence == storage.PostgreSQL:
			q = "SELECT * FROM events WHERE id = $1 LIMIT 1;"
		case constants.Persistence == storage.MySQL:
			q = "SELECT * FROM events WHERE id = ? LIMIT 1;"
		}

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
							"reason":  "Not Found",
							"message": "Event not found",
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

		desc, _ := strconv.ParseBool(c.QueryParam("desc"))

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
					Message: "Unprocessable Entity",
					Errors: []map[string]interface{}{
						{
							"reason":  "Unprocessable Entity",
							"message": "The ID provided cannot be processed as integer",
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
							"reason":  "Internal server error",
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
							"reason":  "Internal Server Error",
							"message": "The ID provided was rejected, not valid",
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
								"reason":  "Conflict",
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

		desc, _ := strconv.ParseBool(c.QueryParam("desc"))

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
							"reason":  "Unprocessable entity",
							"message": "The ID provided cannot be processed as integer",
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
					Code:    400,
					Message: "Bad Request",
					Errors: []map[string]interface{}{
						{
							"reason":  "Bad Request",
							"message": "The ID provided was rejected, not valid",
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
								"reason":  "Conflict",
								"message": "An error was logged while trying to process the payload",
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

		eventId, err := strconv.Atoi(c.Param("event-id"))
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"event_id":       eventId,
					"participant_id": 0,
				},
				Error: models.Error{
					Code:    422,
					Message: "Unprocessable Entity",
					Errors: []map[string]interface{}{
						{
							"reason":  "Unprocessable Entity",
							"message": "The event ID parameter provided cannot be processed as integer",
						},
					},
				},
			})
		}

		participantId, err := strconv.Atoi(c.Param("participant-id"))
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"event_id":       eventId,
					"participant_id": participantId,
				},
				Error: models.Error{
					Code:    422,
					Message: "Unprocessable entity",
					Errors: []map[string]interface{}{
						{
							"reason":  "Unprocessable Entity",
							"message": "The participant ID parameter provided cannot be processed as integer",
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
					"event_id":       eventId,
					"participant_id": participantId,
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
					"event_id":       eventId,
					"participant_id": participantId,
				},
				Error: models.Error{
					Code:    500,
					Message: "Internal server error",
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
		err = stmt.QueryRowContext(ctx, eventId, participantId).Scan(&tview.Id, &tview.Participant, &tview.Event)
		if err != nil {
			return c.JSON(http.StatusNotFound, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"event_id":       eventId,
					"participant_id": participantId,
				},
				Error: models.Error{
					Code:    404,
					Message: "Not Found",
					Errors: []map[string]interface{}{
						{
							"reason":  "Not Found",
							"message": "The event or participant was not found",
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
				"event_id":       eventId,
				"participant_id": participantId,
			},
			Data: tview,
		})
	}
}
