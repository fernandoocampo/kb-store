package kbs

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// KBID defines kb id.
type KBID string

// UserID defines kb id.
type UserID string

// EventID defines kb id.
type EventID string

// OrderByField defines fields you can use to order queries.
type OrderByField string

// NewKB contains data to request the creation of a new kb.
type NewKB struct {
	UserID   UserID  `json:"user_id"`
	UserName string  `json:"username"`
	Content  string  `json:"content"`
	EventID  EventID `json:"event_id"`
}

// UpdateKB contains data to request the update of a new kb.
type UpdateKB struct {
	ID         KBID    `json:"id"`
	UserID     UserID  `json:"user_id"`
	UserName   string  `json:"username"`
	Content    string  `json:"content"`
	EventID    EventID `json:"event_id"`
	UpdateDate int64   `json:"update_date"`
}

// KB contains kb data.
type KB struct {
	ID           KBID    `json:"id"`
	UserID       UserID  `json:"user_id"`
	UserName     string  `json:"username"`
	Content      string  `json:"content"`
	EventID      EventID `json:"event_id"`
	CreationDate int64   `json:"creation_date"`
	UpdateDate   int64   `json:"update_date"`
}

// ValidationError define kb validation logic.
type ValidationError struct {
	KBs []string
}

// QueryFilter contains data for query filters.
type QueryFilter struct {
	EventID     string
	OrderBy     OrderByField
	PageNumber  uint8
	RowsPerPage uint8
}

// GetKBWithIDResult standard roesponse for get a KB with an ID.
type GetKBWithIDResult struct {
	KB  *KB
	Err string
}

// CreateKBResult standard response for create KB.
type CreateKBResult struct {
	ID  KBID
	Err string
}

// UpdateKBResult standard response for updating a kb.
type UpdateKBResult struct {
	Err string
}

// DeleteKBResult standard response for deleting a kb.
type DeleteKBResult struct {
	Err string
}

// SearchKBsResult contains search kbs result data.
type SearchKBsResult struct {
	KBs         []KB
	Total       int
	Page        uint8
	RowsPerPage uint8
}

// SearchKBsDataResult standard roespnse for get a KB with an ID.
type SearchKBsDataResult struct {
	SearchResult SearchKBsResult
	Err          string
}

const (
	// EmptyKBID is the kb id that empty or nil.
	EmptyKBID         = KBID("")
	EmptyOrderByField = OrderByField("")

	PageNumberDefault  = uint8(1)
	RowsPerPageDefault = uint8(10)
)

// order by field possible values
const (
	UserIDField OrderByField = "UserID"
)

func (e *ValidationError) addErrorKB(kb string) {
	e.KBs = append(e.KBs, kb)
}

func (e *ValidationError) Error() string {
	// TODO improve this part
	return fmt.Sprintf("invalid kb data: %+v", e.KBs)
}

func newKBID() KBID {
	return KBID(uuid.New().String())
}

func buildNewKB(newKB NewKB) KB {
	return KB{
		ID:           newKBID(),
		UserID:       newKB.UserID,
		UserName:     newKB.UserName,
		Content:      newKB.Content,
		EventID:      newKB.EventID,
		CreationDate: time.Now().UTC().Unix(),
	}
}

func validKBToUpdate(kb UpdateKB) error {
	err := new(ValidationError)

	if kb.ID == "" {
		err.addErrorKB("kb id cannot be empty")
	}

	if kb.UserID == "" {
		err.addErrorKB("user id cannot be empty")
	}

	if kb.UserName == "" {
		err.addErrorKB("user name cannot be empty")
	}

	if kb.EventID == "" {
		err.addErrorKB("event id cannot be empty")
	}

	if kb.Content == "" {
		err.addErrorKB("content cannot be empty")
	}

	if len(err.KBs) > 0 {
		return err
	}

	return nil
}

func (q QueryFilter) isInvalid() bool {
	return false
}

func (q *QueryFilter) fillDefaultValues() {
	if q.OrderBy == EmptyOrderByField {
		q.OrderBy = UserIDField
	}

	if q.PageNumber == 0 {
		q.PageNumber = PageNumberDefault
	}

	if q.RowsPerPage == 0 {
		q.RowsPerPage = RowsPerPageDefault
	}
}

// newGetKBWithIDResult create a new GetKBWithIDResult
func newGetKBWithIDResult(kb *KB, err error) GetKBWithIDResult {
	var errkb string
	if err != nil {
		errkb = err.Error()
	}
	return GetKBWithIDResult{
		KB:  kb,
		Err: errkb,
	}
}

// newCreateKBResult create a new CreateKBResponse
func newCreateKBResult(id KBID, err error) CreateKBResult {
	var errkb string
	if err != nil {
		errkb = err.Error()
	}
	return CreateKBResult{
		ID:  id,
		Err: errkb,
	}
}

// newUpdateKBResult udpate a new UpdateKBResponse
func newUpdateKBResult(err error) UpdateKBResult {
	var errkb string
	if err != nil {
		errkb = err.Error()
	}
	return UpdateKBResult{
		Err: errkb,
	}
}

// newDeleteKBResult udpate a new DeleteKBResponse
func newDeleteKBResult(err error) DeleteKBResult {
	var errkb string
	if err != nil {
		errkb = err.Error()
	}
	return DeleteKBResult{
		Err: errkb,
	}
}

// newSearchKBsResult create a new SearchKBsResult
func newSearchKBsDataResult(result SearchKBsResult, err error) SearchKBsDataResult {
	var errkb string
	if err != nil {
		errkb = err.Error()
	}
	return SearchKBsDataResult{
		SearchResult: result,
		Err:          errkb,
	}
}

func (m KBID) String() string {
	return string(m)
}

func (u UserID) String() string {
	return string(u)
}

func (e EventID) String() string {
	return string(e)
}

func (u UpdateKB) UpdateDateString() string {
	return strconv.FormatInt(u.UpdateDate, 10)
}

func (u *UpdateKB) fillUpdateTime() {
	u.UpdateDate = time.Now().UTC().Unix()
}
