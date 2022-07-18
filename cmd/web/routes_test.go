package main

import (
	"fmt"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/nanmenkaimak/bookings/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing
	default:
		t.Error(fmt.Sprintf("type is not *chi.Mux, but it is %T", v))
	}
}
