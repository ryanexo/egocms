package format

import `fmt`

func Code(e int64, len int64) string {
    str := fmt.Sprintf("%%0%dd", len)
    return fmt.Sprintf(str, e)
}
