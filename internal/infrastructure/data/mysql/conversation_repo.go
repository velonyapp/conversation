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

var _ repo.Conversation = (*conversationRepo)(nil)

type conversationRepo struct {
	db         *sql.DB
	dispatcher *domainevent.Dispatcher
	log        *slog.Logger
}

func NewConversationRepo(
	db *sql.DB,
	dispatcher *domainevent.Dispatcher,
	log *slog.Logger,
) repo.Conversation {
	return &conversationRepo{
		db:         db,
		dispatcher: dispatcher,
		log:        log,
	}
}

type repoConversationScanner interface {
	Scan(dest ...any) error
}

func (r *conversationRepo) FindByID(ctx context.Context, conversationID vo.ConversationID) (*entity.Conversation, error) {
	const query = `
		SELECT
			id,
			title,
			description,
			create_time,
			update_time,
			version
		FROM conversations
		WHERE id = ?
		LIMIT 1
	`

	row := executor(ctx, r.db).QueryRowContext(ctx, query, conversationID.String())

	conversation, err := scanRepoConversation(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		r.log.ErrorContext(ctx,
			"failed to find conversation",
			slog.String("conversation_id", conversationID.String()),
			slog.Any("error", err),
		)
		return nil, err
	}

	return conversation, nil
}

func (r *conversationRepo) Save(ctx context.Context, conversation *entity.Conversation) error {
	if !conversation.IsDirty() {
		return nil
	}

	if conversation.IsDeleted() {
		if !conversation.IsNew() {
			const query = `
				DELETE FROM conversations
				WHERE id = ?
			`

			_, err := executor(ctx, r.db).ExecContext(ctx, query, conversation.ID().String())
			if err != nil {
				r.log.ErrorContext(ctx,
					"failed to delete conversation",
					slog.String("conversation_id", conversation.ID().String()),
					slog.Any("error", err),
				)
				return err
			}
		}
	} else {
		if conversation.IsNew() {
			const query = `
				INSERT INTO conversations (
					id,
					title,
					description,
					create_time,
					update_time,
					version
				)
				VALUES (?, ?, ?, ?, ?, ?)
			`

			if _, err := executor(ctx, r.db).ExecContext(ctx, query,
				conversation.ID().String(),
				conversation.Title().String(),
				conversation.Description().String(),
				conversation.CreateTime(),
				conversation.UpdateTime(),
				conversation.Version(),
			); err != nil {
				r.log.ErrorContext(ctx,
					"failed to insert conversation",
					slog.String("conversation_id", conversation.ID().String()),
					slog.Any("error", err),
				)
				return err
			}
		} else {
			const query = `
			    UPDATE conversations
			    SET
			        title = ?,
			        description = ?,
			        update_time = ?,
			        version = ?
			    WHERE id = ?
			        AND version = ?
			`

			result, err := executor(ctx, r.db).ExecContext(ctx, query,
				conversation.Title().String(),
				conversation.Description().String(),
				conversation.UpdateTime(),
				conversation.Version(),
				conversation.ID().String(),
				conversation.Version()-1,
			)
			if err != nil {
				r.log.ErrorContext(ctx,
					"failed to update conversation",
					slog.String("conversation_id", conversation.ID().String()),
					slog.Any("error", err),
				)
				return err
			}
			rowsAffected, err := result.RowsAffected()
			if err != nil {
				r.log.ErrorContext(ctx,
					"failed to get conversation rows affected",
					slog.String("conversation_id", conversation.ID().String()),
					slog.Any("error", err),
				)
				return err
			}
			if rowsAffected == 0 {
				return repo.ErrConversationVersionConflict
			}
		}
	}

	for _, domainEvent := range conversation.PullEvents() {
		if err := r.dispatcher.Dispatch(ctx, domainEvent); err != nil {
			r.log.ErrorContext(ctx,
				"failed to dispatch conversation domain event",
				slog.String("conversation_id", conversation.ID().String()),
				slog.Any("error", err),
			)
			return err
		}
	}

	conversation.MarkPersisted()

	return nil
}

func scanRepoConversation(scanner repoConversationScanner) (*entity.Conversation, error) {
	var (
		idRaw          string
		titleRaw       string
		descriptionRaw string
		createTime     time.Time
		updateTime     time.Time
		version        uint32
	)

	if err := scanner.Scan(
		&idRaw,
		&titleRaw,
		&descriptionRaw,
		&createTime,
		&updateTime,
		&version,
	); err != nil {
		return nil, err
	}

	id, _ := vo.NewConversationID(idRaw)
	title, _ := vo.NewConversationTitle(titleRaw)
	description, _ := vo.NewConversationDescription(descriptionRaw)

	return entity.ReconstituteConversation(
		id,
		title,
		description,
		createTime,
		updateTime,
		nil,
		version,
	), nil
}
