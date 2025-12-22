package validator

import `bytes`

type ValidationError struct {
    Field  string `json:"field"`
    Reason string `json:"reason"`
}

type ValidationErrors []ValidationError

func (errs ValidationErrors) Error() string {
    var buf bytes.Buffer
    for _, e := range errs {
        buf.WriteString(e.Field + ":" + e.Reason + "\n")
    }
    return buf.String()
}
