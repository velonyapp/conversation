package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/velonyapp/conversation/internal/domain/entity"
	"github.com/velonyapp/conversation/internal/domain/repo"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

var _ repo.ConversationUser = (*conversationUserRepo)(nil)

type conversationUserRepo struct {
	db *sql.DB
}

func NewConversationUserRepo(db *sql.DB) repo.ConversationUser {
	return &conversationUserRepo{
		db: db,
	}
}

type repoConversationUserScanner interface {
	Scan(dest ...any) error
}

func (r *conversationUserRepo) FindByID(ctx context.Context, conversationUserID vo.ConversationUserID) (*entity.ConversationUser, error) {
	const query = `
		SELECT
			id,
			conversation_id,
			user_id,
			role,
			create_time,
			update_time,
			version
		FROM conversation_users
		WHERE id = ?
		LIMIT 1
	`

	row := executor(ctx, r.db).QueryRowContext(ctx, query, conversationUserID.String())

	conversationUser, err := scanRepoConversationUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return conversationUser, nil
}

func (r *conversationUserRepo) FindByConversationIDAndUserID(ctx context.Context, conversationID vo.ConversationID, userID vo.UserID) (*entity.ConversationUser, error) {
	const query = `
		SELECT
			id,
			conversation_id,
			user_id,
			role,
			create_time,
			update_time,
			version
		FROM conversation_users
		WHERE conversation_id = ? AND user_id = ?
		LIMIT 1
	`

	row := executor(ctx, r.db).QueryRowContext(ctx, query, conversationID.String(), userID.String())

	conversationUser, err := scanRepoConversationUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return conversationUser, nil
}

func (r *conversationUserRepo) Save(ctx context.Context, conversationUser *entity.ConversationUser) error {
	if !conversationUser.IsDirty() {
		return nil
	}

	if conversationUser.IsRemoved() {
		if !conversationUser.IsNew() {
			const query = `
				DELETE FROM conversation_users
				WHERE id = ?
			`

			_, err := executor(ctx, r.db).ExecContext(ctx, query, conversationUser.ID().String())
			if err != nil {
				return err
			}
		}
	} else {
		if conversationUser.IsNew() {
			const query = `
				INSERT INTO conversation_users (
					id,
					conversation_id,
					user_id,
					role,
					create_time,
					update_time,
					version
				)
				VALUES (?, ?, ?, ?, ?, ?, ?)
			`

			if _, err := executor(ctx, r.db).ExecContext(ctx, query,
				conversationUser.ID().String(),
				conversationUser.ConversationID().String(),
				conversationUser.UserID().String(),
				conversationUser.Role().String(),
				conversationUser.CreateTime(),
				conversationUser.UpdateTime(),
				conversationUser.Version(),
			); err != nil {
				return err
			}
		} else {
			const query = `
			    UPDATE conversation_users
			    SET
			        role = ?,
			        update_time = ?,
			        version = ?
			    WHERE id = ?
			        AND version = ?
			`

			result, err := executor(ctx, r.db).ExecContext(ctx, query,
				conversationUser.Role().String(),
				conversationUser.UpdateTime(),
				conversationUser.Version(),
				conversationUser.ID().String(),
				conversationUser.Version()-1,
			)
			if err != nil {
				return err
			}
			rowsAffected, err := result.RowsAffected()
			if err != nil {
				return err
			}
			if rowsAffected == 0 {
				return repo.ErrConversationUserVersionConflict
			}
		}
	}

	conversationUser.MarkPersisted()

	return nil
}

func scanRepoConversationUser(scanner repoConversationUserScanner) (*entity.ConversationUser, error) {
	var (
		idRaw             string
		conversationIDRaw string
		userIDRaw         string
		roleRaw           string
		createTime        time.Time
		updateTime        time.Time
		version           uint32
	)

	if err := scanner.Scan(
		&idRaw,
		&conversationIDRaw,
		&userIDRaw,
		&roleRaw,
		&createTime,
		&updateTime,
		&version,
	); err != nil {
		return nil, err
	}

	id, _ := vo.NewConversationUserID(idRaw)
	conversationID, _ := vo.NewConversationID(conversationIDRaw)
	userID := vo.NewUserID(userIDRaw)
	role, _ := vo.NewConversationUserRole(roleRaw)

	return entity.ReconstituteConversationUser(
		id,
		conversationID,
		userID,
		role,
		createTime,
		updateTime,
		nil,
		version,
	), nil
}
