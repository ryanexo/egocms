package assembler

import (
    `dpcms/internal/app/user/internal/dto`
    "dpcms/internal/infra/persistence/model"
    `dpcms/internal/types`
)

func BuildUserCreateCommand(data *dto.UserCreateParams) *model.User {
    return &model.User{
        Username: data.Username,
        Password: data.Password,
        Email:    data.Email,
        IP:       data.IP,
        Profile:  &model.UserProfile{},
    }
}

func BuildUserDTO(data *model.User) *dto.User {
    return &dto.User{
        Base: types.Base{
            ID:        data.ID,
            CreatedAt: data.CreatedAt,
            UpdatedAt: data.UpdatedAt,
        },
        Username: data.Username,
        Email:    data.Email,
        IP:       data.IP,
        Status:   data.Status,
        RoleID:   data.RoleID,
        RoleName: data.Role.Name,
        Profile: &dto.UserProfile{
            Avatar:      data.Profile.Avatar,
            Nickname:    data.Profile.Nickname,
            Gender:      data.Profile.Gender,
            Description: data.Profile.Description,
            Country:     data.Profile.Country,
            Province:    data.Profile.Province,
            City:        data.Profile.City,
        },
    }
}

func BuildUserListDTO(data []*model.User) []*dto.User {
    users := make([]*dto.User, 0, len(data))
    for _, user := range data {
        users = append(users, BuildUserDTO(user))
    }
    return users
}

func BuildUserProfileModel(data *dto.UserProfile) *model.UserProfile {
    return &model.UserProfile{
        UserID:      data.UserID,
        Avatar:      data.Avatar,
        Nickname:    data.Nickname,
        Gender:      data.Gender,
        Description: data.Description,
        Country:     data.Country,
        Province:    data.Province,
        City:        data.City,
    }
}
