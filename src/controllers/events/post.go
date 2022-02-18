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

func New() echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			db      = storage.Get(constants.Persistence)
			request = new(models.Event)
			err     error
		)

		if err = c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.post",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    400,
					Message: "Bad request",
					Errors: []map[string]interface{}{
						{
							"reason":  "Bad request",
							"message": "The request body data is not valid",
						},
					},
				},
			})
		}

		if (*request == models.Event{}) {
			return c.JSON(http.StatusBadRequest, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.post",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    400,
					Message: "Bad Request",
					Errors: []map[string]interface{}{
						{
							"reason":  "The request body is empty",
							"message": "Bad Request",
						},
					},
				},
			})
		}

		if err = db.Connect(); err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.post",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    500,
					Message: "Internal Server Error",
					Errors: []map[string]interface{}{
						{
							"reason": "Internal Server Error",
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
		switch constants.Persistence {
		case storage.PostgreSQL:
			q = "INSERT INTO events(name) VALUES($1);"
		case storage.MySQL:
			q = "INSERT INTO events(name) VALUES(?);"
		}

		stmt, err := db.PrepareContext(ctx, q)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.post",
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

		r, err := stmt.ExecContext(ctx, request.Name)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.post",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    400,
					Message: "Bad Request",
					Errors: []map[string]interface{}{
						{
							"reason": "Bad Request",
							"message": "The request body data was rejected, not valid",
						},
					},
				},
			})
		}

		if i, _ := r.RowsAffected(); i == 0 {
			return c.JSON(http.StatusBadRequest, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.post",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    400,
					Message: "Bad Request",
					Errors: []map[string]interface{}{
						{
							"reason":  "Bad Request",
							"message": "Are you following any criteria for insertion?",
						},
					},
				},
			})
		}
		return c.JSON(http.StatusCreated, models.SuccessfulResponse{
			APIVersion: constants.APIVersion,
			Method:     "events.post",
			Context:    c.Request().URL.String(),
		})
	}
}

func NewParticipantByIds() echo.HandlerFunc {
	return func(c echo.Context) error {
		var db = storage.Get(constants.Persistence)

		eventId, err := strconv.Atoi(c.Param("event-id"))
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.post",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"event_id":       eventId,
					"participant_id": 0,
				},
				Error: models.Error{
					Code:    422,
					Message: "Unprocessable entity",
					Errors: []map[string]interface{}{
						{
							"reason":  "Unprocessable entity",
							"message": "The event ID provided cannot be processed as integer",
						},
					},
				},
			})
		}

		participantId, err := strconv.Atoi(c.Param("participant-id"))
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.post",
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
							"reason":  err,
							"message": "The participant ID provided cannot be processed as integer",
						},
					},
				},
			})
		}

		if err = db.Connect(); err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.post",
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
		switch constants.Persistence {
		case storage.PostgreSQL:
			q = "SELECT EXISTS (SELECT t.id AS id, CONCAT(p.firstname, ' ',p.lastname) AS participant, e.name AS event FROM tickets AS t INNER JOIN events AS e ON e.id=t.event INNER JOIN participants AS p ON p.id=t.participant WHERE e.id = $1 AND p.id = $2);"
		case storage.MySQL:
			q = "SELECT EXISTS (SELECT t.id AS id, CONCAT(p.firstname, ' ',p.lastname) AS participant, e.name AS event FROM tickets AS t INNER JOIN events AS e ON e.id=t.event INNER JOIN participants AS p ON p.id=t.participant WHERE e.id = ? AND p.id = ?);"
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		stmt, err := db.PrepareContext(ctx, q)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.post",
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
							"reason":  "Internal server error",
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
		
		var exists bool
		stmt.QueryRowContext(ctx, eventId, participantId).Scan(&exists)
		if exists {
			return c.JSON(http.StatusBadRequest, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method: "events.post",
				Context: c.Request().URL.String(),
				Params: map[string]interface{}{
					"event_id":       eventId,
					"participant_id": participantId,
				},
				Error: models.Error{
					Code: 400,
					Message: "Bad Request",
					Errors: []map[string]interface{}{
						{
							"reason": "Bad Request",
							"message": "The participant was already registered for the event previously",
						},
					},
				},
			})
		}

		switch constants.Persistence {
		case storage.PostgreSQL:
			q = "INSERT INTO tickets(participant, event) VALUES($1, $2);"
		case storage.MySQL:
			q = "INSERT INTO tickets(participant, event) VALUES(?, ?);"
		}

		stmt, err = db.PrepareContext(ctx, q)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.post",
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

		r, err := stmt.ExecContext(ctx, participantId, eventId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.post",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"event_id":       eventId,
					"participant_id": participantId,
				},
				Error: models.Error{
					Code:    400,
					Message: "Bad request",
					Errors: []map[string]interface{}{
						{
							"reason":  "Bad request",
							"message": "The request body data was rejected, not valid",
						},
					},
				},
			})
		}

		if i, _ := r.RowsAffected(); i == 0 {
			return c.JSON(http.StatusBadRequest, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.post",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"event_id":       eventId,
					"participant_id": participantId,
				},
				Error: models.Error{
					Code:    400,
					Message: "Bad Request",
					Errors: []map[string]interface{}{
						{
							"reason":  "Bad Request",
							"message": "Are you following any criteria for insertion?",
						},
					},
				},
			})
		}
		return c.JSON(http.StatusCreated, models.SuccessfulResponse{
			APIVersion: constants.APIVersion,
			Method:     "events.post",
			Context:    c.Request().URL.String(),
			Params: map[string]interface{}{
				"event_id":       eventId,
				"participant_id": participantId,
			},
		})
	}
}
