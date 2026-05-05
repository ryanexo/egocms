package domain

import `cms/internal/infra/persistence/datatype`

type Article struct {
    ID         datatype.SafeUint64
    Url        string
    CategoryId datatype.SafeUint64
    Flag       int16
    Title      string
}
