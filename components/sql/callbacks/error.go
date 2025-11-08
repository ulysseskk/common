package callbacks

import (
	"context"
	errors2 "errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"gitlab.ulyssesk.top/common/common/components/sql/metrics"
	"gitlab.ulyssesk.top/common/common/model/errors"
	"gitlab.ulyssesk.top/common/common/model/rest"
	"gitlab.ulyssesk.top/common/common/trace"
	commonContext "gitlab.ulyssesk.top/common/common/util/context"
	"gorm.io/gorm"
)

const (
	ctxKeyNotAllowRecordNotFound = "_record_not_found_not_allowed"
)

func CreateErrorSolveCallback(f func(ctx context.Context, tableName string, err error) error) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		if db.Error == nil {
			return
		}
		tableName := "unknown"
		if db.Statement != nil && db.Statement.Table != "" {
			tableName = db.Statement.Table
		}
		db.Error = f(db.Statement.Context, tableName, db.Error)
	}
}

func ErrorWithStack(ctx context.Context, tableName string, originErr error) error {
	if ctx != nil {
		_, exist := commonContext.GetValue(ctx, ctxKeyNotAllowRecordNotFound)
		if !exist && errors2.Is(originErr, gorm.ErrRecordNotFound) {
			return nil
		}
	}
	var pqErr *pq.Error
	caller := trace.GetNearestCaller(2)
	var err error
	errMsg := ""
	if errors2.As(originErr, &pqErr) {
		errMsg = pqErr.Message
		err = errors.NewError().WithError(originErr).WithCode(rest.DatabaseError).WithMessage(pqErr.Message)
	} else {
		errMsg = originErr.Error()
		if len(errMsg) > 10 {
			errMsg = errMsg[:10]
		}
		err = errors.NewError().WithError(originErr).WithCode(rest.DatabaseError)
	}
	metrics.RecordSQLError(caller, tableName, errMsg)
	return err
}

func RecordNotFoundNotAllowed(ctx context.Context) context.Context {
	return commonContext.WithObject(ctx, ctxKeyNotAllowRecordNotFound, "")
}

func RestErrorWithStack(ctx context.Context, tableName string, originErr error) error {
	if ctx != nil {
		_, exist := commonContext.GetValue(ctx, ctxKeyNotAllowRecordNotFound)
		if !exist && errors2.Is(originErr, gorm.ErrRecordNotFound) {
			return nil
		}
	}
	var pqErr *pgconn.PgError
	caller := trace.GetNearestCaller(3)
	errMsg := ""
	if errors2.As(originErr, &pqErr) {
		errMsg = pqErr.Message
	} else {
		errMsg = originErr.Error()
		if len(errMsg) > 10 {
			errMsg = errMsg[:10]
		}
	}
	metrics.RecordSQLError(caller, tableName, errMsg)
	err := errors.NewError().WithError(originErr).WithCode(rest.DatabaseError)
	if errors2.Is(originErr, gorm.ErrRecordNotFound) {
		return rest.Error{
			Code:        rest.RequestDataNotExisted,
			Message:     "",
			OriginError: err,
		}
	}
	return rest.Error{
		Code:        rest.DatabaseError,
		Message:     "",
		OriginError: err,
	}
}
