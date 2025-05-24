package password

import (
    `golang.org/x/crypto/bcrypt`
)

func Make(raw string) (string, error) {
    r, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
    return string(r), err
}

func Compare(hash string, password string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
