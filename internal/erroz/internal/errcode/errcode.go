package errcode

import `fmt`

type Code int64

func (e Code) String(len int64) string {
    format := fmt.Sprintf("%%0%dd", len)
    return fmt.Sprintf(format, e)
}
