package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/velonyapp/conversation/internal/domain/entity"
	"github.com/velonyapp/conversation/internal/domain/repo"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

var _ repo.User = (*userRepo)(nil)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) repo.User {
	return &userRepo{
		db: db,
	}
}

type repoUserScanner interface {
	Scan(dest ...any) error
}

func (r *userRepo) FindByID(ctx context.Context, userID vo.UserID) (*entity.User, error) {
	const query = `
		SELECT
			id
		FROM users
		WHERE id = ?
		LIMIT 1
	`

	row := executor(ctx, r.db).QueryRowContext(ctx, query, userID.String())

	user, err := scanRepoUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (r *userRepo) Create(ctx context.Context, user *entity.User) error {
	const query = `
		INSERT INTO users (
			id
		)
		VALUES (?)
	`

	if _, err := executor(ctx, r.db).ExecContext(ctx, query,
		user.ID().String(),
	); err != nil {
		return err
	}

	return nil
}

func (r *userRepo) Delete(ctx context.Context, userID vo.UserID) error {
	const query = `
		DELETE FROM users
		WHERE id = ?
	`

	if _, err := executor(ctx, r.db).ExecContext(ctx, query,
		userID.String(),
	); err != nil {
		return err
	}

	return nil
}

func scanRepoUser(scanner repoUserScanner) (*entity.User, error) {
	var (
		idRaw string
	)

	if err := scanner.Scan(
		&idRaw,
	); err != nil {
		return nil, err
	}

	id := vo.NewUserID(idRaw)

	return entity.ReconstituteUser(
		id,
	), nil
}
