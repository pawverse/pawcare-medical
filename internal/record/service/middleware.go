package service

import (
	"context"

	"github.com/pawverse/pawcare-medical/internal/record/domain"
	"github.com/pawverse/pawcare-core/pkg/utils"
	"go.uber.org/zap"
)

type middleware func(IRecordService) IRecordService

type loggingMiddleware struct {
	logger *zap.Logger
	next   IRecordService
}

func newLoggingMiddleware(logger *zap.Logger) middleware {
	return func(next IRecordService) IRecordService {
		return &loggingMiddleware{logger, next}
	}
}

func (mw *loggingMiddleware) Logger(ctx context.Context) *zap.Logger {
	requestId := ctx.Value(utils.RequestIdContextKey).(string)
	return mw.logger.With(zap.String("request_id", requestId))
}

func (mw *loggingMiddleware) GetById(ctx context.Context, id domain.RecordId) (record *domain.Record, err error) {
	defer func() {
		mw.Logger(ctx).
			Info("GetById",
				zap.String("method", "GetById"),
				zap.String("id", string(id)),
				zap.Stringer("record", record),
				zap.Error(err),
			)
	}()

	return mw.next.GetById(ctx, id)
}

func (mw *loggingMiddleware) GetByPetId(ctx context.Context, petId domain.PetId) (records []*domain.Record, err error) {
	defer func() {
		mw.Logger(ctx).
			Info("GetByPetId",
				zap.String("method", "GetByPetId"),
				zap.String("petId", string(petId)),
				zap.Int("#records", len(records)),
				zap.Error(err),
			)
	}()

	return mw.next.GetByPetId(ctx, petId)
}

func (mw *loggingMiddleware) Create(ctx context.Context, petId domain.PetId, recordInfo domain.RecordInfo) (record *domain.Record, err error) {
	defer func() {
		mw.Logger(ctx).
			Info("Create",
				zap.String("method", "Create"),
				zap.String("petId", string(petId)),
				zap.Stringer("recordInfo", recordInfo),
				zap.Error(err),
			)
	}()

	return mw.next.Create(ctx, petId, recordInfo)
}
