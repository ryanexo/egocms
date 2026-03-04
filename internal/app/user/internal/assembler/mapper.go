package assembler

import (
    `cms/internal/app/user/internal/dto`
    "cms/internal/infra/persist/model"
    `cms/internal/util/types`
)

func ToUserCreateCommand(data *dto.UserCreateParams) *model.User {
    return &model.User{
        Username: data.Username,
        Password: data.Password,
        Email:    data.Email,
        IP:       data.IP,
        Profile:  &model.UserProfile{},
    }
}

func ToUserDTO(data *model.User) *dto.User {
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

func ToUserListDTO(data []*model.User) []*dto.User {
    users := make([]*dto.User, 0, len(data))
    for _, user := range data {
        users = append(users, ToUserDTO(user))
    }
    return users
}

func ToUserProfileModel(data *dto.UserProfile) *model.UserProfile {
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
