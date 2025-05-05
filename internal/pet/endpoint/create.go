package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-playground/validator/v10"
	"github.com/pawverse/pawcare-core/pkg/common"
	"github.com/pawverse/pawcare-medical/internal/pet/domain"
	"github.com/pawverse/pawcare-medical/internal/pet/service"
)

type CreateRequest struct {
	Id     string `json:"id" validate:"required"`
	UserId string `json:"user_id" validate:"required"`
}

type GetResponse struct {
	Id     string `json:"id"`
	UserId string `json:"user_id"`

	Err error `json:"-"`
}

func (r GetResponse) Failed() error {
	return r.Err
}

func makeCreateEndpoint(petService service.IPetService) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		req, ok := request.(CreateRequest)
		if !ok {
			return GetResponse{Err: common.ErrCastRequest}, nil
		}

		if err := validator.New().Struct(req); err != nil {
			return GetResponse{Err: err}, nil
		}

		petId := domain.PetId(req.Id)
		userId := domain.UserId(req.UserId)
		pet, err := petService.Create(ctx, petId, userId)
		if err != nil {
			return GetResponse{Err: err}, nil
		}

		return GetResponse{
			Id:     string(pet.Id),
			UserId: string(pet.UserId),
		}, nil
	}
}
