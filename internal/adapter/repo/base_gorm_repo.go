package repo

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rmscoal/moilerplate/internal/app/usecase"
	"gorm.io/gorm"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// SQLSTATE is the error code
type SQLSTATE string

const (
	DuplicateError  SQLSTATE = "23505"
	ForeignKeyError SQLSTATE = "23503"
	// TODO: Add more states
)

func (s SQLSTATE) String() string {
	return string(s)
}

// baseRepo is the base repository
// where all other repo is inherited
// from.
type baseRepo struct {
	db          *gorm.DB
	constraints map[string]string
	tracer      trace.Tracer
}

var gormRepo *baseRepo

func InitBaseRepo(db *gorm.DB, skipRegistry bool) error {
	gormRepo = &baseRepo{
		db:          db,
		constraints: make(map[string]string, 0),
		tracer:      otel.Tracer("gorm_repo_tracer"),
	}

	if !skipRegistry {
		if err := gormRepo.registerIndexes(); err != nil {
			return err
		}

		if err := gormRepo.registerForeignKeys(); err != nil {
			return err
		}
	}

	return nil
}

func (repo *baseRepo) registerIndexes() error {
	rows, err := repo.db.Raw(`
	SELECT
			indexname AS index_name,
			string_agg(replace(attname, '_', ' '), ', and ') AS indexed_columns
	FROM
			pg_indexes
	JOIN
			pg_index ON pg_indexes.indexname::regclass = pg_index.indexrelid
	JOIN
			pg_attribute ON pg_attribute.attrelid = pg_indexes.tablename::regclass
			AND pg_attribute.attnum = ANY(pg_index.indkey)
	WHERE
			schemaname = 'public' -- Change this if your indexes are in a different schema
			AND
			indexname LIKE 'idx_%'
			AND
			attname NOT IN ('deleted_at')
	GROUP BY
			indexname;
	`).Rows()
	if err != nil {
		return err
	}

	for rows.Next() {
		var idx, column string
		rows.Scan(&idx, &column)
		repo.constraints[idx] = column
	}

	return nil
}

func (repo *baseRepo) registerForeignKeys() error {
	rows, err := repo.db.Raw(`
	SELECT
			conname AS foreign_key_name,
			replace(confrelid::regclass::text, '_', ' ') AS referenced_table
	FROM
			pg_constraint
	JOIN
			pg_attribute AS a ON a.attnum = ANY(conkey) AND a.attrelid = conrelid
	JOIN
			pg_attribute AS af ON af.attnum = ANY(confkey) AND af.attrelid = confrelid
	WHERE
			confrelid IS NOT NULL;
	`).Rows()
	if err != nil {
		return err
	}

	for rows.Next() {
		var fkey, table string
		rows.Scan(&fkey, &table)
		repo.constraints[fkey] = table
	}

	return nil
}

func (repo *baseRepo) DetectConstraintError(err error) error {
	if err == nil {
		return nil
	}

	if pgErr, ok := err.(*pgconn.PgError); ok {
		switch SQLSTATE(pgErr.Code) {
		case DuplicateError:
			return fmt.Errorf("%s already exists", repo.constraints[pgErr.ConstraintName])
		case ForeignKeyError:
			return fmt.Errorf("association error to %s", repo.constraints[pgErr.ConstraintName])
		}
	}

	return usecase.ErrUnexpected
}

func (repo *baseRepo) DetectNotFoundError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return usecase.ErrNotFound
	}

	return usecase.ErrUnexpected
}
