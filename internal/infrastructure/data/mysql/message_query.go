package mysql

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/velonyapp/conversation/internal/application/common"
	"github.com/velonyapp/conversation/internal/application/query"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

var _ query.Message = (*messageQuery)(nil)

type messageQuery struct {
	db  *sql.DB
	log *slog.Logger
}

func NewMessageQuery(
	db *sql.DB,
	log *slog.Logger,
) query.Message {
	return &messageQuery{
		db:  db,
		log: log,
	}
}

type queryMessageScanner interface {
	Scan(dest ...any) error
}

func (q *messageQuery) Find(ctx context.Context, messageID vo.MessageID, scope query.MessageScope) (*common.MessageResult, error) {
	const query = `
		SELECT
		    m.id,
		    m.user_id,
		    m.content,
		    m.create_time,
		    m.update_time
		FROM messages m
		INNER JOIN conversation_users cu
		    ON m.conversation_id = cu.conversation_id
		WHERE m.id = ?
		    AND cu.conversation_id = ?
		    AND cu.user_id = ?
		LIMIT 1
	`

	row := executor(ctx, q.db).QueryRowContext(ctx, query,
		messageID.String(),
		scope.ConversationID.String(),
		scope.UserID.String(),
	)

	message, err := scanQueryMessage(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		q.log.ErrorContext(ctx,
			"failed to find message",
			slog.String("message_id", messageID.String()),
			slog.String("conversation_id", scope.ConversationID.String()),
			slog.String("user_id", scope.UserID.String()),
			slog.Any("error", err),
		)
		return nil, err
	}

	return &message, nil
}

func (q *messageQuery) BatchFind(ctx context.Context, messageIDs []vo.MessageID, scope query.MessageScope) (map[vo.MessageID]common.MessageResult, error) {
	results := make(map[vo.MessageID]common.MessageResult, len(messageIDs))
	if len(messageIDs) == 0 {
		return results, nil
	}

	placeholders := make([]string, len(messageIDs))
	args := make([]any, len(messageIDs))

	for i, messageID := range messageIDs {
		placeholders[i] = "?"
		args[i] = messageID.String()
	}

	args = append(args, scope.ConversationID.String())
	args = append(args, scope.UserID.String())

	query := `
		SELECT
			m.id,
			m.user_id,
			m.content,
			m.create_time,
			m.update_time
		FROM messages m
		INNER JOIN conversation_users cu
			ON m.conversation_id = cu.conversation_id
		WHERE m.id IN (` + strings.Join(placeholders, ", ") + `)
			AND cu.conversation_id = ?
			AND cu.user_id = ?
	`

	rows, err := executor(ctx, q.db).QueryContext(ctx, query, args...)
	if err != nil {
		q.log.ErrorContext(ctx,
			"failed to batch find messages",
			slog.Int("message_count", len(messageIDs)),
			slog.String("conversation_id", scope.ConversationID.String()),
			slog.String("user_id", scope.UserID.String()),
			slog.Any("error", err),
		)
		return nil, err
	}
	defer rows.Close()

	messages := make(map[string]common.MessageResult, len(messageIDs))
	for rows.Next() {
		message, err := scanQueryMessage(rows)
		if err != nil {
			q.log.ErrorContext(ctx,
				"failed to scan message during batch find",
				slog.Any("error", err),
			)
			return nil, err
		}

		messages[message.ID] = message
	}

	if err := rows.Err(); err != nil {
		q.log.ErrorContext(ctx,
			"failed to iterate messages during batch find",
			slog.Any("error", err),
		)
		return nil, err
	}

	for _, messageID := range messageIDs {
		if message, ok := messages[messageID.String()]; ok {
			results[messageID] = message
		}
	}

	return results, nil
}

type messageCursor struct {
	Time time.Time `json:"time"`
	ID   string    `json:"id"`
}

func (q *messageQuery) List(ctx context.Context, scope query.MessageScope, pageSize int32, pageToken string) ([]common.MessageResult, string, error) {
	var cursor messageCursor
	if pageToken != "" {
		data, err := base64.RawURLEncoding.DecodeString(pageToken)
		if err != nil {
			q.log.ErrorContext(ctx,
				"failed to decode message page token",
				slog.String("conversation_id", scope.ConversationID.String()),
				slog.String("user_id", scope.UserID.String()),
				slog.Any("error", err),
			)
			return nil, "", err
		}

		if err := json.Unmarshal(data, &cursor); err != nil {
			q.log.ErrorContext(ctx,
				"failed to unmarshal message page token",
				slog.String("conversation_id", scope.ConversationID.String()),
				slog.String("user_id", scope.UserID.String()),
				slog.Any("error", err),
			)
			return nil, "", err
		}
	}

	const query = `
		SELECT
			m.id,
			m.user_id,
			m.content,
			m.create_time,
			m.update_time
		FROM messages m
		INNER JOIN conversation_users cu
			ON m.conversation_id = cu.conversation_id
		WHERE cu.conversation_id = ?
			AND cu.user_id = ?
			AND (? = '' OR (m.create_time, m.id) < (?, ?))
		ORDER BY
			m.create_time DESC,
			m.id DESC
		LIMIT ?
	`

	rows, err := executor(ctx, q.db).QueryContext(ctx, query,
		scope.ConversationID.String(),
		scope.UserID.String(),
		pageToken,
		cursor.Time,
		cursor.ID,
		int64(pageSize)+1,
	)
	if err != nil {
		q.log.ErrorContext(ctx,
			"failed to list messages",
			slog.String("conversation_id", scope.ConversationID.String()),
			slog.String("user_id", scope.UserID.String()),
			slog.Int("page_size", int(pageSize)),
			slog.Any("error", err),
		)
		return nil, "", err
	}
	defer rows.Close()

	messages := make([]common.MessageResult, 0, pageSize)
	for rows.Next() {
		if len(messages) == int(pageSize) {
			nextCursor := messageCursor{
				Time: messages[len(messages)-1].CreateTime,
				ID:   messages[len(messages)-1].ID,
			}

			data, err := json.Marshal(nextCursor)
			if err != nil {
				q.log.ErrorContext(ctx,
					"failed to marshal message page token",
					slog.String("conversation_id", scope.ConversationID.String()),
					slog.String("user_id", scope.UserID.String()),
					slog.Any("error", err),
				)
				return nil, "", err
			}

			return messages, base64.RawURLEncoding.EncodeToString(data), nil
		}

		message, err := scanQueryMessage(rows)
		if err != nil {
			q.log.ErrorContext(ctx,
				"failed to scan message during list",
				slog.Any("error", err),
			)
			return nil, "", err
		}

		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		q.log.ErrorContext(ctx,
			"failed to iterate messages during list",
			slog.Any("error", err),
		)
		return nil, "", err
	}

	return messages, "", nil
}

func scanQueryMessage(scanner queryMessageScanner) (common.MessageResult, error) {
	var (
		id         string
		userIDRaw  sql.NullString
		content    string
		createTime time.Time
		updateTime time.Time
	)

	if err := scanner.Scan(
		&id,
		&userIDRaw,
		&content,
		&createTime,
		&updateTime,
	); err != nil {
		return common.MessageResult{}, err
	}

	var userID *string
	if userIDRaw.Valid {
		value := userIDRaw.String
		userID = &value
	}

	return common.MessageResult{
		ID:         id,
		UserID:     userID,
		Content:    content,
		CreateTime: createTime,
		UpdateTime: updateTime,
	}, nil
}
