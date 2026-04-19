package valueobject

import (
    `fmt`
    `regexp`
)

var regex = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_-]+[a-zA-Z0-9]$")

type URLPath struct {
    path string
}

func NewURLPath(path string) (URLPath, error) {
    if !regex.MatchString(path) {
        return URLPath{}, fmt.Errorf("invalid url path: %s", path)
    }
    return URLPath{path}, nil
}

func (s URLPath) String() string {
    return s.path
}
