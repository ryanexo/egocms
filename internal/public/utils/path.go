package utils

import `regexp`

var regex = regexp.MustCompile("^[a-zA-Z0-9]+[a-zA-Z0-9_-]*[a-zA-Z0-9]*$")

func IsValidPath(p string) bool {
    return regex.MatchString(p)
}
