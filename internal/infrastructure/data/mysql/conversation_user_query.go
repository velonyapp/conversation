package mysql

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/velonyapp/conversation/internal/application/common"
	"github.com/velonyapp/conversation/internal/application/query"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

var _ query.ConversationUser = (*conversationUserQuery)(nil)

type conversationUserQuery struct {
	db *sql.DB
}

func NewConversationUserQuery(db *sql.DB) query.ConversationUser {
	return &conversationUserQuery{
		db: db,
	}
}

type queryConversationUserScanner interface {
	Scan(dest ...any) error
}

func (q *conversationUserQuery) Find(ctx context.Context, conversationUserID vo.ConversationUserID, scope query.ConversationUserScope) (*common.ConversationUserResult, error) {
	const query = `
		SELECT
		    cu.id,
		    cu.user_id,
		    cu.role,
		    cu.create_time,
		    cu.update_time
		FROM conversation_users cu
		JOIN conversation_users me
		    ON me.conversation_id = cu.conversation_id
		WHERE cu.id = ?
		    AND me.conversation_id = ?
		    AND me.user_id = ?
		LIMIT 1
	`

	row := executor(ctx, q.db).QueryRowContext(ctx, query,
		conversationUserID.String(),
		scope.ConversationID.String(),
		scope.UserID.String(),
	)

	conversationUser, err := scanQueryConversationUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &conversationUser, nil
}

func (q *conversationUserQuery) BatchFind(ctx context.Context, conversationUserIDs []vo.ConversationUserID, scope query.ConversationUserScope) (map[vo.ConversationUserID]common.ConversationUserResult, error) {
	results := make(map[vo.ConversationUserID]common.ConversationUserResult, len(conversationUserIDs))
	if len(conversationUserIDs) == 0 {
		return results, nil
	}

	placeholders := make([]string, len(conversationUserIDs))
	args := make([]any, len(conversationUserIDs))

	for i, conversationUserID := range conversationUserIDs {
		placeholders[i] = "?"
		args[i] = conversationUserID.String()
	}

	args = append(args, scope.ConversationID.String())
	args = append(args, scope.UserID.String())

	query := `
	    SELECT
	        cu.id,
	        cu.user_id,
	        cu.role,
	        cu.create_time,
	        cu.update_time
	    FROM conversation_users cu
	    JOIN conversation_users me
	        ON me.conversation_id = cu.conversation_id
	    WHERE cu.id IN (` + strings.Join(placeholders, ", ") + `)
	        AND me.conversation_id = ?
	        AND me.user_id = ?
	`

	rows, err := executor(ctx, q.db).QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	conversationUsers := make(map[string]common.ConversationUserResult, len(conversationUserIDs))
	for rows.Next() {
		conversationUser, err := scanQueryConversationUser(rows)
		if err != nil {
			return nil, err
		}

		conversationUsers[conversationUser.ID] = conversationUser
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, conversationUserID := range conversationUserIDs {
		if conversationUser, ok := conversationUsers[conversationUserID.String()]; ok {
			results[conversationUserID] = conversationUser
		}
	}

	return results, nil
}

type conversationUserCursor struct {
	Time time.Time `json:"time"`
	ID   string    `json:"id"`
}

func (q *conversationUserQuery) List(ctx context.Context, scope query.ConversationUserScope, pageSize int32, pageToken string) ([]common.ConversationUserResult, string, error) {
	var cursor conversationUserCursor
	if pageToken != "" {
		data, err := base64.RawURLEncoding.DecodeString(pageToken)
		if err != nil {
			return nil, "", err
		}

		if err := json.Unmarshal(data, &cursor); err != nil {
			return nil, "", err
		}
	}

	const query = `
	    SELECT
	        cu.id,
	        cu.user_id,
	        cu.role,
	        cu.create_time,
	        cu.update_time
	    FROM conversation_users cu
	    JOIN conversation_users me
	        ON me.conversation_id = cu.conversation_id
	    WHERE me.conversation_id = ?
	        AND me.user_id = ?
	        AND (? = '' OR (cu.create_time, cu.id) < (?, ?))
	    ORDER BY
	        cu.create_time DESC,
	        cu.id DESC
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
		return nil, "", err
	}
	defer rows.Close()

	conversationUsers := make([]common.ConversationUserResult, 0, pageSize)
	for rows.Next() {
		if len(conversationUsers) == int(pageSize) {
			nextCursor := conversationUserCursor{
				Time: conversationUsers[len(conversationUsers)-1].CreateTime,
				ID:   conversationUsers[len(conversationUsers)-1].ID,
			}

			data, err := json.Marshal(nextCursor)
			if err != nil {
				return nil, "", err
			}

			return conversationUsers, base64.RawURLEncoding.EncodeToString(data), nil
		}

		conversationUser, err := scanQueryConversationUser(rows)
		if err != nil {
			return nil, "", err
		}

		conversationUsers = append(conversationUsers, conversationUser)
	}
	if err := rows.Err(); err != nil {
		return nil, "", err
	}

	return conversationUsers, "", nil
}

func scanQueryConversationUser(scanner queryConversationUserScanner) (common.ConversationUserResult, error) {
	var (
		id         string
		userID     string
		role       string
		createTime time.Time
		updateTime time.Time
	)

	if err := scanner.Scan(
		&id,
		&userID,
		&role,
		&createTime,
		&updateTime,
	); err != nil {
		return common.ConversationUserResult{}, err
	}

	return common.ConversationUserResult{
		ID:         id,
		UserID:     userID,
		Role:       role,
		CreateTime: createTime,
		UpdateTime: updateTime,
	}, nil
}
