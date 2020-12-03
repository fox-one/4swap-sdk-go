package fswap

type (
	Pagination struct {
		NextCursor string `json:"next_cursor,omitempty"`
		HasNext    bool   `json:"has_next,omitempty"`
	}
)
