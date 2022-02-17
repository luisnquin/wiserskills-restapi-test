package participants

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

		desc, _ := strconv.ParseBool(c.QueryParam("desc"))

		if err = db.Connect(); err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "participants.get",
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

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		var q string
		if desc {
			q = "SELECT * FROM participants ORDER BY id DESC;"
		} else {
			q = "SELECT * FROM participants;"
		}

		stmt, err := db.PrepareContext(ctx, q)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "participants.get",
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
				Method:     "participants.get",
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

		var participants models.Participants
		for rows.Next() {
			var p models.Participant
			if err = rows.Scan(&p.Id, &p.Firstname, &p.Lastname, &p.Age); err != nil {
				return c.JSON(http.StatusConflict, models.BadResponse{
					APIVersion: constants.APIVersion,
					Method:     "participants.get",
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
			participants = append(participants, p)
		}

		if len(participants) == 0 {
			return c.JSON(http.StatusNoContent, models.SuccessfulResponse{
				APIVersion: constants.APIVersion,
				Method:     "participants.get",
				Context:    c.Request().URL.String(),
			})
		}

		return c.JSON(http.StatusOK, models.SuccessfulResponse{
			APIVersion: constants.APIVersion,
			Method:     "participants.get",
			Context:    c.Request().URL.String(),
			Data:       participants,
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
				Method:     "participants.get",
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
				Method:     "participants.get",
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

		q := "SELECT * FROM participants WHERE id = ? LIMIT 1;"
		if constants.Persistence == storage.PostgreSQL {
			q = sqlx.Rebind(sqlx.DOLLAR, q)
		}

		stmt, err := db.PrepareContext(ctx, q)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "participants.get",
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

		var p models.Participant
		err = stmt.QueryRowContext(ctx, id).Scan(&p.Id, &p.Firstname, &p.Lastname, &p.Age)
		if err != nil {
			return c.JSON(http.StatusNotFound, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "participants.get",
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
			Method:     "participants.get",
			Context:    c.Request().URL.String(),
			Params: map[string]interface{}{
				"id": id,
			},
			Data: p,
		})
	}
}

func FetchTicketsById() echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			db  = storage.Get(constants.Persistence)
			err error
		)

		desc, _ := strconv.ParseBool(c.QueryParam("desc"))

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "participants.get",
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
				Method:     "participants.get",
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
			q = "SELECT t.id AS id, CONCAT(p.firstname, ' ',p.lastname) AS participant, e.name AS event FROM tickets AS t INNER JOIN events AS e ON e.id=t.event INNER JOIN participants AS p ON p.id=t.participant WHERE p.id = $1 ORDER BY t.id DESC;"
		case constants.Persistence == storage.MySQL && desc:
			q = "SELECT t.id AS id, CONCAT(p.firstname, ' ',p.lastname) AS participant, e.name AS event FROM tickets AS t INNER JOIN events AS e ON e.id=t.event INNER JOIN participants AS p ON p.id=t.participant WHERE p.id = ? ORDER BY t.id DESC;"
		case constants.Persistence == storage.PostgreSQL:
			q = "SELECT t.id AS id, CONCAT(p.firstname, ' ',p.lastname) AS participant, e.name AS event FROM tickets AS t INNER JOIN events AS e ON e.id=t.event INNER JOIN participants AS p ON p.id=t.participant WHERE p.id = $1 ORDER BY t.id ASC;"
		case constants.Persistence == storage.MySQL:
			q = "SELECT t.id AS id, CONCAT(p.firstname, ' ',p.lastname) AS participant, e.name AS event FROM tickets AS t INNER JOIN events AS e ON e.id=t.event INNER JOIN participants AS p ON p.id=t.participant WHERE p.id = ? ORDER BY t.id ASC;"
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		stmt, err := db.PrepareContext(ctx, q)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "participants.get",
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

		rows, err := stmt.QueryContext(ctx, id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "participants.get",
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

		var tviews models.TicketViews

		for rows.Next() {
			var tview models.TicketView

			if err = rows.Scan(&tview.Id, &tview.Participant, &tview.Event); err != nil {
				return c.JSON(http.StatusConflict, models.BadResponse{
					APIVersion: constants.APIVersion,
					Method:     "participants.get",
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
				Method:     "participants.get",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
			})
		}
		return c.JSON(http.StatusOK, models.SuccessfulResponse{
			APIVersion: constants.APIVersion,
			Method:     "participants.get",
			Context:    c.Request().URL.String(),
			Params: map[string]interface{}{
				"id": id,
			},
			Data: tviews,
		})
	}
}


