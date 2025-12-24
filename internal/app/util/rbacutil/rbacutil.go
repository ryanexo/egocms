package rbachelper

import (
    `strconv`
    `strings`
)

const (
    subjectPrefix = "ROLE::"
)

func GetRoleSubject(id uint64) string {
    return subjectPrefix + strconv.FormatUint(id, 10)
}

func ParseRoleSubject(names ...string) ([]uint64, error) {
    idList := make([]uint64, 0, len(names))
    for _, inheritRoleName := range names {
        id, err := strconv.ParseUint(strings.Replace(inheritRoleName, subjectPrefix, "", 1), 10, 64)
        if err != nil {
            return nil, err
        }
        idList = append(idList, id)
    }
    return idList, nil
}
