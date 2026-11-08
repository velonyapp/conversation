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

var _ repo.Message = (*messageRepo)(nil)

type messageRepo struct {
	db         *sql.DB
	dispatcher *domainevent.Dispatcher
}

func NewMessageRepo(
	db *sql.DB,
	dispatcher *domainevent.Dispatcher,
) repo.Message {
	return &messageRepo{
		db:         db,
		dispatcher: dispatcher,
	}
}

type repoMessageScanner interface {
	Scan(dest ...any) error
}

func (r *messageRepo) FindByID(ctx context.Context, messageID vo.MessageID) (*entity.Message, error) {
	const query = `
		SELECT
			id,
			conversation_id,
			user_id,
			content,
			create_time,
			update_time,
			version
		FROM messages
		WHERE id = ?
		LIMIT 1
	`

	row := executor(ctx, r.db).QueryRowContext(ctx, query, messageID.String())

	message, err := scanRepoMessage(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return message, nil
}

func (r *messageRepo) FindLastByConversationID(ctx context.Context, conversationID vo.ConversationID) (*entity.Message, error) {
	const query = `
		SELECT
			m.id,
			m.conversation_id,
			m.user_id,
			m.content,
			m.create_time,
			m.update_time,
			m.version
		FROM messages m
		INNER JOIN conversations c
			ON m.id = c.last_message_id
		WHERE c.id = ?
		LIMIT 1
	`

	row := executor(ctx, r.db).QueryRowContext(ctx, query, conversationID.String())

	message, err := scanRepoMessage(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return message, nil
}

func (r *messageRepo) Save(ctx context.Context, message *entity.Message) error {
	if !message.IsDirty() {
		return nil
	}

	if message.IsDeleted() {
		if !message.IsNew() {
			const conversationQuery = `
        	    UPDATE conversations c
        	    LEFT JOIN (
        	        SELECT
        	            id,
        	            create_time
        	        FROM messages
        	        WHERE conversation_id = ?
        	            AND id <> ?
        	        ORDER BY
        	            create_time DESC,
        	            id DESC
        	        LIMIT 1
        	    ) lm ON TRUE
        	    SET
        	        c.last_message_id = lm.id,
        	        c.last_message_create_time = lm.create_time
        	    WHERE c.id = ?
        	        AND c.last_message_id = ?
        	`

			if _, err := executor(ctx, r.db).ExecContext(ctx, conversationQuery,
				message.ConversationID().String(),
				message.ID().String(),
				message.ConversationID().String(),
				message.ID().String(),
			); err != nil {
				return err
			}

			const query = `
        	    DELETE FROM messages
        	    WHERE id = ?
        	`

			if _, err := executor(ctx, r.db).ExecContext(ctx, query,
				message.ID().String(),
			); err != nil {
				return err
			}
		}
	} else {
		var userID any
		if message.HasUser() {
			userID = message.UserID().String()
		}

		if message.IsNew() {
			const query = `
				INSERT INTO messages (
					id,
					conversation_id,
					user_id,
					content,
					create_time,
					update_time,
					version
				)
				VALUES (?, ?, ?, ?, ?, ?, ?)
			`

			if _, err := executor(ctx, r.db).ExecContext(ctx, query,
				message.ID().String(),
				message.ConversationID().String(),
				userID,
				message.Content().String(),
				message.CreateTime(),
				message.UpdateTime(),
				message.Version(),
			); err != nil {
				return err
			}

			const conversationQuery = `
				UPDATE conversations
				SET
					last_message_id = ?,
					last_message_create_time = ?
				WHERE id = ?
					AND (last_message_id IS NULL OR (last_message_create_time, last_message_id) < (?, ?))
			`

			if _, err := executor(ctx, r.db).ExecContext(ctx, conversationQuery,
				message.ID().String(),
				message.CreateTime(),
				message.ConversationID().String(),
				message.CreateTime(),
				message.ID().String(),
			); err != nil {
				return err
			}
		} else {
			const query = `
			    UPDATE messages
			    SET
			        content = ?,
			        update_time = ?,
			        version = ?
			    WHERE id = ?
			        AND version = ?
			`

			result, err := executor(ctx, r.db).ExecContext(ctx, query,
				message.Content().String(),
				message.UpdateTime(),
				message.Version(),
				message.ID().String(),
				message.Version()-1,
			)
			if err != nil {
				return err
			}
			rowsAffected, err := result.RowsAffected()
			if err != nil {
				return err
			}
			if rowsAffected == 0 {
				return repo.ErrMessageVersionConflict
			}
		}
	}

	for _, domainEvent := range message.PullEvents() {
		if err := r.dispatcher.Dispatch(ctx, domainEvent); err != nil {
			return err
		}
	}

	message.MarkPersisted()

	return nil
}

func scanRepoMessage(scanner repoMessageScanner) (*entity.Message, error) {
	var (
		idRaw             string
		conversationIDRaw string
		userIDRaw         sql.NullString
		contentRaw        string
		createTime        time.Time
		updateTime        time.Time
		version           uint32
	)

	if err := scanner.Scan(
		&idRaw,
		&conversationIDRaw,
		&userIDRaw,
		&contentRaw,
		&createTime,
		&updateTime,
		&version,
	); err != nil {
		return nil, err
	}

	id, _ := vo.NewMessageID(idRaw)
	conversationID, _ := vo.NewConversationID(conversationIDRaw)
	var userID *vo.UserID
	if userIDRaw.Valid {
		value := vo.NewUserID(userIDRaw.String)
		userID = &value
	}
	content, _ := vo.NewMessageContent(contentRaw)

	return entity.ReconstituteMessage(
		id,
		conversationID,
		userID,
		content,
		createTime,
		updateTime,
		nil,
		version,
	), nil
}
