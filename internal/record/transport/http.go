package transport

import (
	"github.com/pawverse/pawcare-medical/internal/record/endpoint"
	"github.com/pawverse/pawcare-core/pkg/common"
	"github.com/pawverse/pawcare-core/pkg/utils"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func RegisterHTTPRoutes(r *mux.Router, enpoints endpoint.Set, logger *zap.Logger) {
	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(common.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(kithttp.DefaultErrorEncoder),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerBefore(utils.RequestIdHTTPToContext()),
	}
	options = append(options, common.HTTPLoggingServerOptions(logger)...)

	createHandler := kithttp.NewServer(
		enpoints.CreateEndpoint,
		common.DecodeJSONRequest[endpoint.CreateRequest],
		common.EncodeJSONResponse,
		options...,
	)

	getByPetIdHandler := kithttp.NewServer(
		enpoints.GetByPetIdEndpoint,
		common.DecodePathParameters[endpoint.GetByIdRequest],
		common.EncodeJSONResponse,
		options...,
	)

	getByIdHandler := kithttp.NewServer(
		enpoints.GetByIdEndpoint,
		common.DecodePathParameters[endpoint.GetByIdRequest],
		common.EncodeJSONResponse,
		options...,
	)

	r.Methods("GET").Path("/pets/{id}/records").Handler(getByPetIdHandler)
	r.Methods("POST").Path("/pets/{id}/records").Handler(createHandler)
	r.Methods("GET").Path("/records/{id}").Handler(getByIdHandler)
}
