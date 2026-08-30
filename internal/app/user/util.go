package user

import (
	"cms/internal/app/user/api"
	"cms/internal/app/user/model"
	"cms/internal/public/apitype"
	"cms/internal/public/jsontype"
	"cms/internal/public/utils"
)

func buildUser(params api.UserCreateParams) (*model.User, error) {
	user := &model.User{IP: &params.IP, Profile: &model.UserProfile{}}
	if err := user.SetUsername(params.Username); err != nil {
		return nil, err
	}
	if err := user.SetEmail(params.Email); err != nil {
		return nil, err
	}
	passwordHash, err := utils.Password(params.Password).Generate()
	if err != nil {
		return nil, err
	}
	if err = user.SetPasswordHash(passwordHash); err != nil {
		return nil, err
	}
	return user, nil
}

func buildProfile(params api.UserProfileParams) (*model.UserProfile, error) {
	profile := &model.UserProfile{UserID: params.UserID.Uint64()}
	setters := []func() error{
		func() error { return profile.SetAvatar(params.Avatar) },
		func() error { return profile.SetNickname(params.Nickname) },
		func() error { return profile.SetGender(params.Gender) },
		func() error { return profile.SetDescription(params.Description) },
		func() error { return profile.SetCountry(params.Country) },
		func() error { return profile.SetProvince(params.Province) },
		func() error { return profile.SetCity(params.City) },
	}
	for _, set := range setters {
		if err := set(); err != nil {
			return nil, err
		}
	}
	return profile, nil
}

func toUserAPI(user *model.User) *api.User {
	result := &api.User{
		Base: apitype.Base{
			ID:        jsontype.SafeUint64(user.ID),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Username:   user.Username,
		Email:      user.Email,
		VerifiedAt: user.VerifiedAt,
		IP:         user.IP,
		Status:     user.Status,
		RoleName:   user.RoleName,
	}
	if user.RoleID != nil {
		roleID := jsontype.SafeUint64(*user.RoleID)
		result.RoleID = &roleID
	}
	if user.Profile != nil {
		result.Profile = &api.UserProfile{
			Avatar:      user.Profile.Avatar,
			Nickname:    user.Profile.Nickname,
			Gender:      user.Profile.Gender,
			Description: user.Profile.Description,
			Country:     user.Profile.Country,
			Province:    user.Profile.Province,
			City:        user.Profile.City,
		}
	}
	return result
}

func toUserAPIList(users []*model.User) []*api.User {
	result := make([]*api.User, 0, len(users))
	for _, user := range users {
		result = append(result, toUserAPI(user))
	}
	return result
}
