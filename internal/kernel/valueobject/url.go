package valueobject

import (
    `errors`
    `fmt`
    `net/url`
)

var ErrInvalidURLScheme = errors.New("invalid url scheme")

type URL struct {
    u *url.URL
}

func NewURL(value string) (URL, error) {
    u, err := url.ParseRequestURI(value)
    if err != nil {
        return URL{}, err
    }
    if u.Scheme != "http" && u.Scheme != "https" {
        return URL{}, fmt.Errorf("%w: %s", ErrInvalidURLScheme, u.Scheme)
    }
    return URL{u}, nil
}

func (s URL) String() string {
    return s.u.String()
}
