package transport

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pawverse/pawcare-medical/internal/record/endpoint"
	"github.com/pawverse/pawcare-medical/internal/record/transport"
	"go.uber.org/zap"
)

func MakeHTTPServer(endpoints endpoint.Set, logger *zap.Logger) http.Handler {
	router := mux.NewRouter()
	apiGroup := router.PathPrefix("/api/v1").Subrouter()

	transport.RegisterHTTPRoutes(apiGroup, endpoints, logger)

	return router
}
