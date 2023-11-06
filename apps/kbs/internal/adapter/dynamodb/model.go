package dynamodb

import "github.com/fernandoocampo/kb-store/apps/kbs/internal/kbs"

type KB struct {
	ID           string `json:"id" dynamodbav:"id"`
	UserID       string `json:"user_id" dynamodbav:"user_id"`
	UserName     string `json:"username" dynamodbav:"username"`
	Content      string `json:"content" dynamodbav:"content"`
	EventID      string `json:"event_id" dynamodbav:"event_id"`
	CreationDate int64  `json:"creation_date" dynamodbav:"creation_date"`
	UpdateDate   int64  `json:"update_date" dynamodbav:"update_date"`
}

// transformKB transforms new kb to a repository kb.
func (u KB) toRepositoryKB() kbs.KB {
	return kbs.KB{
		ID:           kbs.KBID(u.ID),
		UserID:       kbs.UserID(u.UserID),
		UserName:     u.UserName,
		Content:      u.Content,
		EventID:      kbs.EventID(u.EventID),
		CreationDate: u.CreationDate,
		UpdateDate:   u.UpdateDate,
	}
}

// transformKB transforms new kb to a kb.
func transformKB(kb kbs.KB) KB {
	return KB{
		ID:           kb.ID.String(),
		UserID:       kb.UserID.String(),
		UserName:     kb.UserName,
		Content:      kb.Content,
		EventID:      kb.EventID.String(),
		CreationDate: kb.CreationDate,
		UpdateDate:   kb.UpdateDate,
	}
}
