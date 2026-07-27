package stringx

import (
    `unicode/utf8`
)

func TruncateUTF8(s string, n int) string {
    if n <= 0 {
        return ""
    }
    if n >= len(s) {
        return s
    }
    i := 0
    for i < n {
        _, size := utf8.DecodeRuneInString(s[i:])
        i += size
    }
    return s[:i]
}
