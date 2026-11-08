package mysql

import (
	"context"
	"database/sql"
	"errors"
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
}

func NewConversationRepo(
	db *sql.DB,
	dispatcher *domainevent.Dispatcher,
) repo.Conversation {
	return &conversationRepo{
		db:         db,
		dispatcher: dispatcher,
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
				return err
			}
			rowsAffected, err := result.RowsAffected()
			if err != nil {
				return err
			}
			if rowsAffected == 0 {
				return repo.ErrConversationVersionConflict
			}
		}
	}

	for _, domainEvent := range conversation.PullEvents() {
		if err := r.dispatcher.Dispatch(ctx, domainEvent); err != nil {
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
