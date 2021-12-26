package form

// DeleteParams delete body params
type DeleteParams struct {
	IDs []uint `json:"ids" binding:"required"`
}

// ValueParams ...
type ValueParams struct {
	Value string `json:"value" binding:"required"`
}

// SortColumn sort hiih column ner bolon asc, des
type SortColumn struct {
	Field string `json:"field"`
	Order string `json:"order"`
}
