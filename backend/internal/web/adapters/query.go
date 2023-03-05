package adapters

import (
	"net/http"

	"github.com/hay-kot/homebox/backend/pkgs/server"
)

// Query is a server.Handler that decodes a query from the request and calls the provided function.
func Query[T any, Y any](f AdapterFunc[T, Y], ok int) server.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		q, err := decodeQuery[T](r)
		if err != nil {
			return err
		}

		res, err := f(r.Context(), q)
		if err != nil {
			return err
		}

		return server.Respond(w, ok, res)
	}
}

// QueryID is a server.Handler that decodes a query and an ID from the request and calls the provided function.
func QueryID[T any, Y any](param string, f IDFunc[T, Y], ok int) server.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ID, err := routeUUID(r, param)
		if err != nil {
			return err
		}

		q, err := decodeQuery[T](r)
		if err != nil {
			return err
		}

		res, err := f(r.Context(), ID, q)
		if err != nil {
			return err
		}

		return server.Respond(w, ok, res)
	}
}