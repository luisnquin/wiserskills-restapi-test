package events

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/luisnquin/restapi-technical-test/src/constants"
	"github.com/luisnquin/restapi-technical-test/src/models"
	"github.com/luisnquin/restapi-technical-test/src/storage"
)

func RemoveById() echo.HandlerFunc {
	return func(c echo.Context) error {
		var db = storage.Get(constants.Persistence)

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.delete",
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
							"message": "The ID parameter cannot be used as integer",
						},
					},
				},
			})
		}

		if err = db.Connect(); err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.delete",
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
		switch constants.Persistence {
		case storage.PostgreSQL:
			q = "DELETE FROM events WHERE id = $1;"
		case storage.MySQL:
			q = "DELETE FROM events WHERE id = ?;"
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		stmt, err := db.PrepareContext(ctx, q)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.delete",
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

		r, err := stmt.ExecContext(ctx, id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.delete",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    400,
					Message: "Bad request",
					Errors: []map[string]interface{}{
						{
							"reason":  "Bad Request",
							"message": "The provided ID was rejected, not valid",
						},
					},
				},
			})
		}
		if i, _ := r.RowsAffected(); i == 0 {
			return c.JSON(http.StatusNotFound, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.delete",
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
							"message": "The event was not found",
						},
					},
				},
			})
		}
		return c.JSON(http.StatusOK, models.SuccessfulResponse{
			APIVersion: constants.APIVersion,
			Method:     "events.delete",
			Context:    c.Request().URL.String(),
			Params: map[string]interface{}{
				"id": id,
			},
		})
	}
}

func RemoveByIdWithParticipants() echo.HandlerFunc {
	return func(c echo.Context) error {
		var db = storage.Get(constants.Persistence)

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.delete",
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
							"message": "The ID parameter cannot be used as integer",
						},
					},
				},
			})
		}

		if err = db.Connect(); err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.delete",
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
		switch constants.Persistence {
		case storage.PostgreSQL:
			q = "DELETE FROM participants WHERE id IN (SELECT participant FROM tickets WHERE event = $1);"
		case storage.MySQL:
			q = "DELETE FROM participants WHERE id IN (SELECT participant FROM tickets WHERE event = ?);"
		}

		stmt, err := db.PrepareContext(ctx, q)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.delete",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    400,
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

		_, err = stmt.ExecContext(ctx, id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.delete",
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

		switch constants.Persistence {
		case storage.PostgreSQL:
			q = "DELETE FROM events WHERE id = $1;"
		case storage.MySQL:
			q = "DELETE FROM events WHERE id = ?;"
		}

		stmt, err = db.PrepareContext(ctx, q)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.delete",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    400,
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

		r, err := stmt.ExecContext(ctx, id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.delete",
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
		if i, _ := r.RowsAffected(); i == 0 {
			return c.JSON(http.StatusNotFound, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "events.delete",
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
							"message": fmt.Sprintf("The event with ID %d was not found", id),
						},
					},
				},
			})
		}
		return c.JSON(http.StatusOK, models.SuccessfulResponse{
			APIVersion: constants.APIVersion,
			Method:     "events.delete",
			Context:    c.Request().URL.String(),
			Params: map[string]interface{}{
				"id": id,
			},
		})
	}
}
