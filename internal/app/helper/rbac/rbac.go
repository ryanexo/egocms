package rbac

import (
    `strconv`
    `strings`
)

const (
    subjectPrefix = "ROLE::"
)

func GetRoleSubject(id int64) string {
    return subjectPrefix + strconv.FormatInt(id, 10)
}

func ParseRoleSubject(names ...string) ([]int64, error) {
    idList := make([]int64, 0, len(names))
    for _, inheritRoleName := range names {
        id, err := strconv.ParseInt(strings.Replace(inheritRoleName, subjectPrefix, "", 1), 10, 64)
        if err != nil {
            return nil, err
        }
        idList = append(idList, id)
    }
    return idList, nil
}
