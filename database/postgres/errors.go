package postgres

import (
	stderrs "errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rlapenok/toolbox/errors"
)

// MapError converts database/pgx errors into toolbox errors with appropriate codes and reasons
func MapError(err error) *errors.Error {
	if err == nil {
		return nil
	}

	// if it's already a toolbox error, pass it through
	if tbErr, ok := err.(*errors.Error); ok {
		return tbErr
	}

	var pgErr *pgconn.PgError
	if stderrs.As(err, &pgErr) {
		// common PostgreSQL error codes: https://www.postgresql.org/docs/current/errcodes-appendix.html
		switch pgErr.Code {
		case "23505": // unique_violation
			return errors.New(errors.Conflict, "unique constraint violation").
				WithReason(errors.ReasonConflict).
				WithDetails(map[string]any{
					"code":         pgErr.Code,
					"constraint":   pgErr.ConstraintName,
					"table":        pgErr.TableName,
					"schema":       pgErr.SchemaName,
					"column":       pgErr.ColumnName,
					"detail":       pgErr.Detail,
					"message":      pgErr.Message,
					"where":        pgErr.Where,
					"severity":     pgErr.Severity,
					"routine":      pgErr.Routine,
					"internal_pos": pgErr.InternalPosition,
				})

		case "23503": // foreign_key_violation
			return errors.New(errors.Conflict, "foreign key violation").
				WithReason(errors.ReasonConflict).
				WithDetails(map[string]any{
					"code":       pgErr.Code,
					"constraint": pgErr.ConstraintName,
					"table":      pgErr.TableName,
					"schema":     pgErr.SchemaName,
					"detail":     pgErr.Detail,
					"message":    pgErr.Message,
				})

		case "23514": // check_violation
			return errors.New(errors.BadRequest, "check constraint violation").
				WithReason(errors.ReasonBadRequest).
				WithDetails(map[string]any{
					"code":       pgErr.Code,
					"constraint": pgErr.ConstraintName,
					"table":      pgErr.TableName,
					"schema":     pgErr.SchemaName,
					"detail":     pgErr.Detail,
					"message":    pgErr.Message,
				})

		case "23502": // not_null_violation
			return errors.New(errors.BadRequest, "null value in column violates not-null constraint").
				WithReason(errors.ReasonBadRequest).
				WithDetails(map[string]any{
					"code":    pgErr.Code,
					"table":   pgErr.TableName,
					"schema":  pgErr.SchemaName,
					"column":  pgErr.ColumnName,
					"detail":  pgErr.Detail,
					"message": pgErr.Message,
				})

		case "22P02": // invalid_text_representation
			return errors.New(errors.BadRequest, "invalid text representation").
				WithReason(errors.ReasonBadRequest).
				WithDetails(map[string]any{
					"code":    pgErr.Code,
					"detail":  pgErr.Detail,
					"message": pgErr.Message,
				})

		case "40001": // serialization_failure
			return errors.New(errors.Conflict, "serialization failure").
				WithReason(errors.ReasonConflict).
				WithDetails(map[string]any{
					"code":    pgErr.Code,
					"detail":  pgErr.Detail,
					"message": pgErr.Message,
				})

		case "40P01": // deadlock_detected
			return errors.New(errors.Conflict, "deadlock detected").
				WithReason(errors.ReasonConflict).
				WithDetails(map[string]any{
					"code":    pgErr.Code,
					"detail":  pgErr.Detail,
					"message": pgErr.Message,
				})
		}

		// default mapping for unhandled PG errors
		return errors.New(errors.Internal, pgErr.Message).
			WithReason(errors.ReasonInternal).
			WithDetails(map[string]any{
				"code":    pgErr.Code,
				"detail":  pgErr.Detail,
				"message": pgErr.Message,
				"where":   pgErr.Where,
			})
	}

	// non-PG error -> Internal
	return errors.New(errors.Internal, err.Error()).
		WithReason(errors.ReasonInternal)
}
