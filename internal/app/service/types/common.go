package types

type QueryByIdParam struct {
    ID int64 `validate:"required" json:"id"`
}
