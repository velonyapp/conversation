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

var _ query.Conversation = (*conversationQuery)(nil)

type conversationQuery struct {
	db  *sql.DB
	log *slog.Logger
}

func NewConversationQuery(
	db *sql.DB,
	log *slog.Logger,
) query.Conversation {
	return &conversationQuery{
		db:  db,
		log: log,
	}
}

type queryConversationScanner interface {
	Scan(dest ...any) error
}

func (q *conversationQuery) Find(ctx context.Context, conversationID vo.ConversationID, scope query.ConversationScope) (*common.ConversationResult, error) {
	const query = `
		SELECT
			c.id,
			c.title,
			c.description,
			c.last_message_id,
			c.create_time,
			c.update_time
		FROM conversations c
		INNER JOIN conversation_users cu
			ON c.id = cu.conversation_id
		WHERE cu.conversation_id = ?
			AND cu.user_id = ?
		LIMIT 1
	`

	row := executor(ctx, q.db).QueryRowContext(ctx, query,
		conversationID.String(),
		scope.UserID.String(),
	)

	conversation, err := scanQueryConversation(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		q.log.ErrorContext(ctx,
			"failed to find conversation",
			slog.String("conversation_id", conversationID.String()),
			slog.String("user_id", scope.UserID.String()),
			slog.Any("error", err),
		)
		return nil, err
	}

	return &conversation, nil
}

func (q *conversationQuery) BatchFind(ctx context.Context, conversationIDs []vo.ConversationID, scope query.ConversationScope) (map[vo.ConversationID]common.ConversationResult, error) {
	results := make(map[vo.ConversationID]common.ConversationResult, len(conversationIDs))
	if len(conversationIDs) == 0 {
		return results, nil
	}

	placeholders := make([]string, len(conversationIDs))
	args := make([]any, len(conversationIDs))

	for i, conversationID := range conversationIDs {
		placeholders[i] = "?"
		args[i] = conversationID.String()
	}

	args = append(args, scope.UserID.String())

	query := `
		SELECT
			c.id,
			c.title,
			c.description,
			c.last_message_id,
			c.create_time,
			c.update_time
		FROM conversations c
		INNER JOIN conversation_users cu
			ON c.id = cu.conversation_id
		WHERE cu.conversation_id IN (` + strings.Join(placeholders, ", ") + `)
			AND cu.user_id = ?
	`

	rows, err := executor(ctx, q.db).QueryContext(ctx, query, args...)
	if err != nil {
		q.log.ErrorContext(ctx,
			"failed to batch find conversations",
			slog.Int("conversation_count", len(conversationIDs)),
			slog.String("user_id", scope.UserID.String()),
			slog.Any("error", err),
		)
		return nil, err
	}
	defer rows.Close()

	conversations := make(map[string]common.ConversationResult, len(conversationIDs))
	for rows.Next() {
		conversation, err := scanQueryConversation(rows)
		if err != nil {
			q.log.ErrorContext(ctx,
				"failed to scan conversation during batch find",
				slog.Any("error", err),
			)
			return nil, err
		}

		conversations[conversation.ID] = conversation
	}

	if err := rows.Err(); err != nil {
		q.log.ErrorContext(ctx,
			"failed to iterate conversations during batch find",
			slog.Any("error", err),
		)
		return nil, err
	}

	for _, conversationID := range conversationIDs {
		if conversation, ok := conversations[conversationID.String()]; ok {
			results[conversationID] = conversation
		}
	}

	return results, nil
}

type conversationCursor struct {
	Time time.Time `json:"time"`
	ID   string    `json:"id"`
}

func (q *conversationQuery) List(ctx context.Context, scope query.ConversationScope, pageSize int32, pageToken string) ([]common.ConversationResult, string, error) {
	var cursor conversationCursor
	if pageToken != "" {
		data, err := base64.RawURLEncoding.DecodeString(pageToken)
		if err != nil {
			q.log.ErrorContext(ctx,
				"failed to decode conversation page token",
				slog.String("user_id", scope.UserID.String()),
				slog.Any("error", err),
			)
			return nil, "", err
		}

		if err := json.Unmarshal(data, &cursor); err != nil {
			q.log.ErrorContext(ctx,
				"failed to unmarshal conversation page token",
				slog.String("user_id", scope.UserID.String()),
				slog.Any("error", err),
			)
			return nil, "", err
		}
	}

	const query = `
		SELECT
			c.id,
			c.title,
			c.description,
			c.last_message_id,
			c.create_time,
			c.update_time,
			c.last_message_create_time
		FROM conversations c
		INNER JOIN conversation_users cu
			ON c.id = cu.conversation_id
		WHERE cu.user_id = ?
			AND (? = '' OR (COALESCE(c.last_message_create_time, c.update_time), c.id) < (?, ?))
		ORDER BY
			COALESCE(c.last_message_create_time, c.update_time) DESC,
			c.id DESC
		LIMIT ?
	`

	rows, err := executor(ctx, q.db).QueryContext(ctx, query,
		scope.UserID.String(),
		pageToken,
		cursor.Time,
		cursor.ID,
		int64(pageSize)+1,
	)
	if err != nil {
		q.log.ErrorContext(ctx,
			"failed to list conversations",
			slog.String("user_id", scope.UserID.String()),
			slog.Int("page_size", int(pageSize)),
			slog.Any("error", err),
		)
		return nil, "", err
	}
	defer rows.Close()

	conversations := make([]common.ConversationResult, 0, pageSize)
	var nextCursor conversationCursor

	for rows.Next() {
		if len(conversations) == int(pageSize) {
			data, err := json.Marshal(nextCursor)
			if err != nil {
				q.log.ErrorContext(ctx,
					"failed to marshal conversation page token",
					slog.String("user_id", scope.UserID.String()),
					slog.Any("error", err),
				)
				return nil, "", err
			}

			return conversations, base64.RawURLEncoding.EncodeToString(data), nil
		}

		var lastMessageCreateTime sql.NullTime
		conversation, err := scanQueryConversation(rows, &lastMessageCreateTime)
		if err != nil {
			q.log.ErrorContext(ctx,
				"failed to scan conversation during list",
				slog.Any("error", err),
			)
			return nil, "", err
		}

		nextCursor.ID = conversation.ID
		nextCursor.Time = conversation.UpdateTime
		if lastMessageCreateTime.Valid {
			nextCursor.Time = lastMessageCreateTime.Time
		}

		conversations = append(conversations, conversation)
	}
	if err := rows.Err(); err != nil {
		q.log.ErrorContext(ctx,
			"failed to iterate conversations during list",
			slog.Any("error", err),
		)
		return nil, "", err
	}

	return conversations, "", nil
}

func scanQueryConversation(scanner queryConversationScanner, lastMessageCreateTime ...*sql.NullTime) (common.ConversationResult, error) {
	var (
		id               string
		title            string
		description      string
		lastMessageIDRaw sql.NullString
		createTime       time.Time
		updateTime       time.Time
	)

	dest := []any{
		&id,
		&title,
		&description,
		&lastMessageIDRaw,
		&createTime,
		&updateTime,
	}
	if len(lastMessageCreateTime) > 0 {
		dest = append(dest, lastMessageCreateTime[0])
	}

	if err := scanner.Scan(dest...); err != nil {
		return common.ConversationResult{}, err
	}

	var lastMessageID *string
	if lastMessageIDRaw.Valid {
		value := lastMessageIDRaw.String
		lastMessageID = &value
	}

	return common.ConversationResult{
		ID:            id,
		Title:         title,
		Description:   description,
		LastMessageID: lastMessageID,
		CreateTime:    createTime,
		UpdateTime:    updateTime,
	}, nil
}
