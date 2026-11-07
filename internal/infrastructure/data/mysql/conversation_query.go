package mysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/velonyapp/conversation/internal/application/common"
	"github.com/velonyapp/conversation/internal/application/query"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

var _ query.Conversation = (*conversationQuery)(nil)

type conversationQuery struct {
	db *sql.DB
}

func NewConversationQuery(db *sql.DB) query.Conversation {
	return &conversationQuery{
		db: db,
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
		return nil, err
	}
	defer rows.Close()

	conversations := make(map[string]common.ConversationResult, len(conversationIDs))
	for rows.Next() {
		conversation, err := scanQueryConversation(rows)
		if err != nil {
			return nil, err
		}

		conversations[conversation.ID] = conversation
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, conversationID := range conversationIDs {
		if conversation, ok := conversations[conversationID.String()]; ok {
			results[conversationID] = conversation
		}
	}

	return results, nil
}

func (q *conversationQuery) List(ctx context.Context, scope query.ConversationScope, pageSize int32, pageToken string) ([]common.ConversationResult, string, error) {
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
		WHERE cu.user_id = ?
			AND (c.xxxxxxx, c.id) < (?, ?)
		ORDER BY
			c.xxxxxxx DESC,
			c.id DESC
		LIMIT ?
	`

	rows, err := executor(ctx, q.db).QueryContext(ctx, query,
		scope.UserID.String(),
		"TODO: c.xxxxxxx",
		"TODO: c.id",
		pageSize,
	)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	conversations := make([]common.ConversationResult, 0)
	for rows.Next() {
		conversation, err := scanQueryConversation(rows)
		if err != nil {
			return nil, "", err
		}

		conversations = append(conversations, conversation)
	}

	return conversations, "TODO: nextPageToken", nil
}

func scanQueryConversation(scanner queryConversationScanner) (common.ConversationResult, error) {
	var (
		id               string
		title            string
		description      string
		lastMessageIDRaw sql.NullString
		createTime       time.Time
		updateTime       time.Time
	)

	if err := scanner.Scan(
		&id,
		&title,
		&description,
		&lastMessageIDRaw,
		&createTime,
		&updateTime,
	); err != nil {
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
