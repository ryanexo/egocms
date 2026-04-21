package valueobject

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

type HashedPassword struct {
    passwd string
}

func NewHashedPassword(passwd string) (HashedPassword, error) {
    r, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
    return HashedPassword{passwd: string(r)}, err
}

func (p HashedPassword) Compare(password string) bool {
    return bcrypt.CompareHashAndPassword([]byte(p.passwd), []byte(password)) == nil
}
