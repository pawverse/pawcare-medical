package endpoint

import (
	"github.com/pawverse/pawcare-medical/internal/pet/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/spf13/viper"
)

type Set struct {
	CreateEndpoint endpoint.Endpoint
}

func NewSet(viper viper.Viper, petService service.IPetService) (Set, error) {
	var createEndpoint endpoint.Endpoint
	{
		createEndpoint = makeCreateEndpoint(petService)
	}

	return Set{
		CreateEndpoint: createEndpoint,
	}, nil
}
