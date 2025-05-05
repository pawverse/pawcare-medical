package service

import (
	"context"

	"github.com/pawverse/pawcare-medical/internal/pet/domain"
	"go.uber.org/zap"
)

type middleware func(IPetService) IPetService

type loggingMiddleware struct {
	logger *zap.Logger
	next   IPetService
}

func newLoggingMiddleware(logger *zap.Logger) middleware {
	return func(next IPetService) IPetService {
		return &loggingMiddleware{logger, next}
	}
}

func (mw *loggingMiddleware) Create(ctx context.Context, id domain.PetId, userId domain.UserId) (pet *domain.Pet, err error) {
	defer func() {
		mw.logger.
			Info("Create",
				zap.String("method", "Create"),
				zap.String("id", string(id)),
				zap.String("userId", string(userId)),
				zap.Stringer("pet", pet),
				zap.Error(err),
			)
	}()
	return mw.next.Create(ctx, id, userId)
}

func (mw *loggingMiddleware) GetById(ctx context.Context, id domain.PetId) (pet *domain.Pet, err error) {
	defer func() {
		mw.logger.
			Info("GetById",
				zap.String("method", "GetById"),
				zap.String("id", string(id)),
				zap.Stringer("pet", pet),
				zap.Error(err),
			)
	}()
	return mw.next.GetById(ctx, id)
}
