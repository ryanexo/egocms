package srvparams

type QueryByResourceID struct {
    ID uint64 `validate:"required" json:"id"`
}
