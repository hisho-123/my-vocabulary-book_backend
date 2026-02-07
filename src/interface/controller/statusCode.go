package controller

import (
	"backend/src/domain"
	"log"
	"net/http"
)

func statusCode(err error) int {
	switch err.Error() {
	case domain.BadRequest:
		return http.StatusBadRequest
	case domain.Unauthorized:
		return http.StatusUnauthorized
	case domain.Forbidden:
		return http.StatusForbidden
	case domain.NotFound:
		return http.StatusNotFound
	case domain.Conflict:
		return http.StatusConflict
	case domain.UnprocessableEntity:
		return http.StatusUnprocessableEntity
	case domain.InternalServerError:
		return http.StatusInternalServerError
	default:
		log.Printf("Unknown error code: %s", err.Error())
		return http.StatusInternalServerError
	}
}
