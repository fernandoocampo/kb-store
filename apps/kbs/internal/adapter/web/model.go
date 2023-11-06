package web

import "github.com/fernandoocampo/kb-store/apps/kbs/internal/kbs"

// Result standard result for the service
type Result struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Errors  []string    `json:"errors"`
}

// KB contains kb data.
type KB struct {
	ID string `json:"id"`
	// UserID kb's name.
	UserID       string `json:"user_id"`
	UserName     string `json:"username"`
	Content      string `json:"content"`
	EventID      string `json:"event_id"`
	CreationDate int64  `json:"creation_date"`
	UpdateDate   int64  `json:"update_date"`
}

// NewKB contains the expected data for a new kb.
type NewKB struct {
	UserID   string `json:"user_id"`
	UserName string `json:"username"`
	Content  string `json:"content"`
	EventID  string `json:"event_id"`
}

// UpdateKB contains the expected data to update an kb.
type UpdateKB struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	UserName string `json:"username"`
	Content  string `json:"content"`
	EventID  string `json:"event_id"`
}

// CreateKBResponse standard response for create KB
type CreateKBResponse struct {
	ID  string `json:"id"`
	Err string `json:"err,omitempty"`
}

// GetKBWithIDResponse standard response for get a KB with an ID.
type GetKBWithIDResponse struct {
	KB  *KB    `json:"kb"`
	Err string `json:"err,omitempty"`
}

// SearchKBsResponse standard response for searching kbs with filters.
type SearchKBsResponse struct {
	KBs *SearchKBsResult `json:"result"`
	Err string           `json:"err,omitempty"`
}

// SearchKBFilter contains filters to search kbs
type SearchKBFilter struct {
	// EventID kb's name.
	EventID string
	// Order by field
	OrderBy string
	// Page page to query
	Page uint8
	// rows per page
	PageSize uint8
}

// SearchKBsResult contains search kbs result data.
type SearchKBsResult struct {
	KBs      []KB  `json:"kbs"`
	Total    int   `json:"total"`
	Page     uint8 `json:"page"`
	PageSize uint8 `json:"page_size"`
}

// toKB transforms new kb to a kb object.
func toKB(kb *kbs.KB) *KB {
	if kb == nil {
		return nil
	}
	webKB := KB{
		ID:           kb.ID.String(),
		UserID:       kb.UserID.String(),
		UserName:     kb.UserName,
		Content:      kb.Content,
		EventID:      kb.EventID.String(),
		CreationDate: kb.CreationDate,
		UpdateDate:   kb.UpdateDate,
	}
	return &webKB
}

// toSearchKBResult transforms new kb to a kb object.
func toSearchKBResult(result *kbs.SearchKBsResult) *SearchKBsResult {
	if result == nil {
		return nil
	}
	kbsFound := make([]KB, 0)
	for _, v := range result.KBs {
		kbFound := toKB(&v)
		kbsFound = append(kbsFound, *kbFound)
	}
	webKB := SearchKBsResult{
		KBs:      kbsFound,
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.RowsPerPage,
	}
	return &webKB
}

func (r Result) NotSuccess() bool {
	return !r.Success
}

func (r Result) ThereAreErrors() bool {
	return len(r.Errors) > 0
}

func (r Result) Failed() bool {
	return (r.NotSuccess() && r.ThereAreErrors())
}

// toKB transforms new kb to a kb object.
func (n *NewKB) toKB() *kbs.NewKB {
	if n == nil {
		return nil
	}
	kbDomain := kbs.NewKB{
		UserID:   kbs.UserID(n.UserID),
		UserName: n.UserName,
		Content:  n.Content,
		EventID:  kbs.EventID(n.EventID),
	}
	return &kbDomain
}

// toKB transforms udpate kb to a kb object.
func (u *UpdateKB) toKB() *kbs.UpdateKB {
	if u == nil {
		return nil
	}
	kbDomain := kbs.UpdateKB{
		ID:       kbs.KBID(u.ID),
		UserID:   kbs.UserID(u.UserID),
		UserName: u.UserName,
		Content:  u.Content,
		EventID:  kbs.EventID(u.EventID),
	}
	return &kbDomain
}

func toCreateKBResponse(kbResult kbs.CreateKBResult) Result {
	var kb Result
	if kbResult.Err == "" {
		kb.Success = true
		kb.Data = kbResult.ID
	}
	if kbResult.Err != "" {
		kb.Errors = []string{kbResult.Err}
	}
	return kb
}

func toUpdateKBResponse(kbResult kbs.UpdateKBResult) Result {
	var kb Result
	if kbResult.Err == "" {
		kb.Success = true
	}
	if kbResult.Err != "" {
		kb.Errors = []string{kbResult.Err}
	}
	return kb
}

func toDeleteKBResponse(kbResult kbs.DeleteKBResult) Result {
	var kb Result
	if kbResult.Err == "" {
		kb.Success = true
	}
	if kbResult.Err != "" {
		kb.Errors = []string{kbResult.Err}
	}
	return kb
}

func toGetKBWithIDResponse(kbResult kbs.GetKBWithIDResult) Result {
	var kb Result
	newKB := toKB(kbResult.KB)
	if kbResult.Err == "" {
		kb.Success = true
		kb.Data = newKB
	}
	if kbResult.Err != "" {
		kb.Errors = []string{kbResult.Err}
	}
	return kb
}

func toSearchKBsResponse(kbResult kbs.SearchKBsDataResult) Result {
	var kb Result

	if kbResult.Err == "" {
		kb.Success = true
		kb.Data = toSearchKBResult(&kbResult.SearchResult)
	}
	if kbResult.Err != "" {
		kb.Errors = []string{kbResult.Err}
	}
	return kb
}

func (s SearchKBFilter) toSearchKBFilter() kbs.QueryFilter {
	return kbs.QueryFilter{
		EventID:     s.EventID,
		PageNumber:  s.Page,
		RowsPerPage: s.PageSize,
		OrderBy:     kbs.OrderByField(s.OrderBy),
	}
}
