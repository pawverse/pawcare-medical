package endpoint

import (
	"github.com/MicahParks/keyfunc/v3"
	"github.com/pawverse/pawcare-medical/internal/record/service"
	"github.com/pawverse/pawcare-core/pkg/common"
	"github.com/go-kit/kit/endpoint"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type Set struct {
	CreateEndpoint     endpoint.Endpoint
	GetByPetIdEndpoint endpoint.Endpoint
	GetByIdEndpoint    endpoint.Endpoint
}

func NewSet(viper viper.Viper, recordService service.IRecordService) (Set, error) {
	kf, err := keyfunc.NewDefault([]string{viper.GetString(common.CertsEndpointKey)})
	if err != nil {
		return Set{}, err
	}

	var createEndpoint endpoint.Endpoint
	{
		createEndpoint = makeCreateEndpoint(recordService)
		createEndpoint = common.NewParser(kf.Keyfunc, jwt.SigningMethodRS256, common.RegisteredClaimsFactory)(createEndpoint)
	}

	var getByPetIdEndpoint endpoint.Endpoint
	{
		getByPetIdEndpoint = makeGetByPetIdEndpoint(recordService)
		getByPetIdEndpoint = common.NewParser(kf.Keyfunc, jwt.SigningMethodRS256, common.RegisteredClaimsFactory)(getByPetIdEndpoint)
	}

	var getByIdEndpoint endpoint.Endpoint
	{
		getByIdEndpoint = makeGetByIdEndpoint(recordService)
		getByIdEndpoint = common.NewParser(kf.Keyfunc, jwt.SigningMethodRS256, common.RegisteredClaimsFactory)(getByIdEndpoint)
	}

	return Set{
		CreateEndpoint:     createEndpoint,
		GetByPetIdEndpoint: getByPetIdEndpoint,
		GetByIdEndpoint:    getByIdEndpoint,
	}, nil
}
