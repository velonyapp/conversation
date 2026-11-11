package mysql

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/velonyapp/conversation/internal/application/domainevent"
	"github.com/velonyapp/conversation/internal/domain/entity"
	"github.com/velonyapp/conversation/internal/domain/repo"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

var _ repo.ConversationUser = (*conversationUserRepo)(nil)

type conversationUserRepo struct {
	db         *sql.DB
	dispatcher *domainevent.Dispatcher
	log        *slog.Logger
}

func NewConversationUserRepo(
	db *sql.DB,
	dispatcher *domainevent.Dispatcher,
	log *slog.Logger,
) repo.ConversationUser {
	return &conversationUserRepo{
		db:         db,
		dispatcher: dispatcher,
		log:        log,
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

		r.log.ErrorContext(ctx,
			"failed to find conversation user",
			slog.String("conversation_user_id", conversationUserID.String()),
			slog.Any("error", err),
		)
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

		r.log.ErrorContext(ctx,
			"failed to find conversation user by conversation ID and user ID",
			slog.String("conversation_id", conversationID.String()),
			slog.String("user_id", userID.String()),
			slog.Any("error", err),
		)
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
				r.log.ErrorContext(ctx,
					"failed to delete conversation user",
					slog.String("conversation_user_id", conversationUser.ID().String()),
					slog.Any("error", err),
				)
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
				r.log.ErrorContext(ctx,
					"failed to insert conversation user",
					slog.String("conversation_user_id", conversationUser.ID().String()),
					slog.Any("error", err),
				)
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
				r.log.ErrorContext(ctx,
					"failed to update conversation user",
					slog.String("conversation_user_id", conversationUser.ID().String()),
					slog.Any("error", err),
				)
				return err
			}
			rowsAffected, err := result.RowsAffected()
			if err != nil {
				r.log.ErrorContext(ctx,
					"failed to get conversation user rows affected",
					slog.String("conversation_user_id", conversationUser.ID().String()),
					slog.Any("error", err),
				)
				return err
			}
			if rowsAffected == 0 {
				return repo.ErrConversationUserVersionConflict
			}
		}
	}

	for _, domainEvent := range conversationUser.PullEvents() {
		if err := r.dispatcher.Dispatch(ctx, domainEvent); err != nil {
			r.log.ErrorContext(ctx,
				"failed to dispatch conversation user domain event",
				slog.String("conversation_user_id", conversationUser.ID().String()),
				slog.Any("error", err),
			)
			return err
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
