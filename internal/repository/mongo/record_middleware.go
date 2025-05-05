package mongo

import (
	"context"

	"github.com/pawverse/pawcare-medical/internal/record/domain"
	"github.com/pawverse/pawcare-core/pkg/utils"
	"go.uber.org/zap"
)

type recordMiddleware func(domain.IRecordRepository) domain.IRecordRepository

type loggingMiddleware struct {
	logger *zap.Logger
	next   domain.IRecordRepository
}

func newLoggingMiddleware(logger *zap.Logger) recordMiddleware {
	return func(next domain.IRecordRepository) domain.IRecordRepository {
		return &loggingMiddleware{logger, next}
	}
}

func (mw *loggingMiddleware) Logger(ctx context.Context) *zap.Logger {
	requestId := ctx.Value(utils.RequestIdContextKey).(string)
	return mw.logger.With(zap.String("request_id", requestId))
}

func (mw *loggingMiddleware) FindById(ctx context.Context, id domain.RecordId) (record *domain.Record, err error) {
	defer func() {
		mw.Logger(ctx).
			Info("FindById",
				zap.String("method", "FindById"),
				zap.String("id", string(id)),
				zap.Stringer("record", record),
				zap.Error(err),
			)
	}()
	return mw.next.FindById(ctx, id)
}

func (mw *loggingMiddleware) FindByPetId(ctx context.Context, petId domain.PetId) (records []*domain.Record, err error) {
	defer func() {
		mw.Logger(ctx).
			Info("FindByPetId",
				zap.String("method", "FindByPetId"),
				zap.String("petId", string(petId)),
				zap.Int("#records", len(records)),
				zap.Error(err),
			)
	}()
	return mw.next.FindByPetId(ctx, petId)
}

func (mw *loggingMiddleware) Create(ctx context.Context, record *domain.Record) (err error) {
	defer func() {
		mw.Logger(ctx).
			Info("Create",
				zap.String("method", "Create"),
				zap.Stringer("record", record),
				zap.Error(err),
			)
	}()
	return mw.next.Create(ctx, record)
}

func (mw *loggingMiddleware) Update(ctx context.Context, record *domain.Record) (err error) {
	defer func() {
		mw.Logger(ctx).
			Info("Update",
				zap.String("method", "Update"),
				zap.Stringer("record", record),
				zap.Error(err),
			)
	}()
	return mw.next.Update(ctx, record)
}
