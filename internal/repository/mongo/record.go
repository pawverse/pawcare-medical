package mongo

import (
	"context"
	"errors"

	"github.com/pawverse/pawcare-medical/internal/record/domain"
	"github.com/pawverse/pawcare-core/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

const (
	RecordCollection = "records"
)

type record struct {
	Id          primitive.ObjectID `bson:"_id"`
	PetId       primitive.ObjectID `bson:"pet_id"`
	Type        string             `bson:"type"`
	Description string             `bson:"description"`
	Date        primitive.DateTime `bson:"date"`
}

func (r *record) ToModel() *domain.Record {
	return &domain.Record{
		Id:         domain.RecordId(r.Id.Hex()),
		PetId:      domain.PetId(r.PetId.Hex()),
		RecordInfo: domain.NewRecordInfo(domain.RecordType(r.Type), r.Description, r.Date.Time()),
	}
}

func toRecordModel(r record) *domain.Record {
	return r.ToModel()
}

func fromRecordModel(r domain.Record) (*record, error) {
	id, err := primitive.ObjectIDFromHex(string(r.Id))
	if err != nil {
		if !errors.Is(err, primitive.ErrInvalidHex) {
			return nil, err
		}
		id = primitive.NewObjectID()
	}

	petId, err := primitive.ObjectIDFromHex(string(r.PetId))
	if err != nil {
		return nil, err
	}

	return &record{
		Id:          id,
		PetId:       petId,
		Type:        string(r.RecordInfo.Type),
		Description: r.RecordInfo.Description,
		Date:        primitive.NewDateTimeFromTime(r.RecordInfo.Date),
	}, nil
}

type recordRepository struct {
	db *mongo.Database
}

func NewRecordRepository(db *mongo.Database, logger *zap.Logger) domain.IRecordRepository {
	var repository domain.IRecordRepository
	repository = &recordRepository{db}
	repository = newLoggingMiddleware(logger)(repository)

	return repository
}

func (r *recordRepository) Collection() *mongo.Collection {
	return r.db.Collection(RecordCollection)
}

func (r *recordRepository) FindById(ctx context.Context, id domain.RecordId) (*domain.Record, error) {
	objectId, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return nil, err
	}

	var result record
	if err := r.Collection().FindOne(ctx, bson.M{"_id": objectId}).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrRecordNotFound
		}

		return nil, err
	}

	return result.ToModel(), nil
}

func (r *recordRepository) FindByPetId(ctx context.Context, petId domain.PetId) ([]*domain.Record, error) {
	objectId, err := primitive.ObjectIDFromHex(string(petId))
	if err != nil {
		return nil, err
	}

	cursor, err := r.Collection().Find(ctx, bson.M{"pet_id": objectId})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	var result []record
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return utils.Map(result, toRecordModel), nil
}

func (r *recordRepository) Create(ctx context.Context, record *domain.Record) error {
	entity, err := fromRecordModel(*record)
	if err != nil {
		return err
	}

	result, err := r.Collection().InsertOne(ctx, entity)
	if err != nil {
		return err
	}

	record.Id = domain.RecordId(result.InsertedID.(primitive.ObjectID).Hex())
	return nil
}

func (r *recordRepository) Update(ctx context.Context, record *domain.Record) error {
	entity, err := fromRecordModel(*record)
	if err != nil {
		return err
	}

	if _, err := r.Collection().UpdateOne(ctx, bson.M{"_id": entity.Id}, bson.M{"$set": entity}); err != nil {
		return err
	}

	return nil
}
