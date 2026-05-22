package password

import (
    `golang.org/x/crypto/bcrypt`
)

type Password string

func (p Password) Generate() (string, error) {
    r, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
    return string(r), err
}

func (p Password) Compare(password string) bool {
    return bcrypt.CompareHashAndPassword([]byte(p), []byte(password)) == nil
}
