package types

// Common parameters related to filtering and such that are common to all get-like operations.
type CommonGetParams struct {
	OrderBy   string
	OrderType string
	Limit     int
	Offset    int
	Filter    string
}
