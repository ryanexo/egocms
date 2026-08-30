package role

import (
    `strconv`
    
    `cms/internal/app/role/model`
    `cms/internal/public/jsontype`
)

const rolePrefix = "role_"

func casbinRoleName(id uint64) string {
    return rolePrefix + strconv.FormatUint(id, 10)
}

func casbinRoleNames(ids []uint64) []string {
    roles := make([]string, 0, len(ids))
    for _, id := range ids {
        roles = append(roles, casbinRoleName(id))
    }
    return roles
}

func buildRole(name, description string) (*model.Role, error) {
    role := &model.Role{}
    if err := role.SetName(name); err != nil {
        return nil, err
    }
    if err := role.SetDescription(description); err != nil {
        return nil, err
    }
    return role, nil
}

func uniqueRoleParams(ids []jsontype.SafeUint64) []uint64 {
    result := make([]uint64, 0, len(ids))
    seen := make(map[uint64]struct{}, len(ids))
    for _, id := range ids {
        value := id.Uint64()
        if _, exists := seen[value]; exists {
            continue
        }
        seen[value] = struct{}{}
        result = append(result, value)
    }
    return result
}
