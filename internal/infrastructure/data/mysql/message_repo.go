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

var _ repo.Message = (*messageRepo)(nil)

type messageRepo struct {
	db *sql.DB
}

func NewMessageRepo(db *sql.DB) repo.Message {
	return &messageRepo{
		db: db,
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
		WHERE conversations.id = ?
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
			const query = `
				DELETE FROM messages
				WHERE id = ?
			`

			_, err := executor(ctx, r.db).ExecContext(ctx, query, message.ID().String())
			if err != nil {
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
