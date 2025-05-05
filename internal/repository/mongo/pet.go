package mongo

import (
	"context"

	"github.com/pawverse/pawcare-medical/internal/pet/domain"
	"github.com/pawverse/pawcare-core/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	PetCollection = "pets"
)

type pet struct {
	Id     primitive.ObjectID `bson:"_id"`
	UserId string             `bson:"user_id"`
}

func (p *pet) ToModel() *domain.Pet {
	return &domain.Pet{
		Id:     domain.PetId(p.Id.Hex()),
		UserId: domain.UserId(p.UserId),
	}
}

func toPetModel(p pet) *domain.Pet {
	return p.ToModel()
}

func fromPetModel(p domain.Pet) (*pet, error) {
	id, err := primitive.ObjectIDFromHex(string(p.Id))
	if err != nil {
		return nil, err
	}

	return &pet{
		Id:     id,
		UserId: string(p.UserId),
	}, nil
}

type petRepository struct {
	db *mongo.Database
}

func NewPetRepository(db *mongo.Database) domain.IPetRepository {
	return &petRepository{db}
}

func (r *petRepository) Collection() *mongo.Collection {
	return r.db.Collection(PetCollection)
}

func (r *petRepository) FindById(ctx context.Context, id domain.PetId) (*domain.Pet, error) {
	objectId, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return nil, err
	}

	var result pet
	if err := r.Collection().FindOne(ctx, bson.M{"_id": objectId}).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrPetNotFound
		}
	}

	return result.ToModel(), nil
}

func (r *petRepository) FindByUserId(ctx context.Context, userId domain.UserId) ([]*domain.Pet, error) {
	cursor, err := r.Collection().Find(ctx, bson.M{"user_id": string(userId)})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrPetNotFound
		}

		return nil, err
	}

	var result []pet
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return utils.Map(result, toPetModel), nil
}

func (r *petRepository) Create(ctx context.Context, pet *domain.Pet) error {
	entity, err := fromPetModel(*pet)
	if err != nil {
		return err
	}

	if _, err := r.Collection().InsertOne(ctx, entity); err != nil {
		return err
	}

	return nil
}
