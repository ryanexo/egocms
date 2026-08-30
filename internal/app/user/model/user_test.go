package model

import (
	"strings"
	"testing"

	"cms/internal/app/user/errno"

	"github.com/stretchr/testify/require"
)

func TestUserSetUsernameValidatesBeforeMutation(t *testing.T) {
	user := &User{Username: "original"}

	err := user.SetUsername(strings.Repeat("a", UsernameMaxLength+1))

	require.ErrorIs(t, err, errno.ErrInvalidUsername)
	require.Equal(t, "original", user.Username)
}

func TestUserSetUsernameAcceptsDatabaseLimit(t *testing.T) {
	username := strings.Repeat("a", UsernameMaxLength)
	user := &User{}

	err := user.SetUsername(username)

	require.NoError(t, err)
	require.Equal(t, username, user.Username)
}

func TestUserSetEmailValidatesBeforeMutation(t *testing.T) {
	user := &User{Email: "original@example.com"}

	err := user.SetEmail("invalid")

	require.ErrorIs(t, err, errno.ErrInvalidEmail)
	require.Equal(t, "original@example.com", user.Email)
}

func TestUserProfileSetNicknameValidatesBeforeMutation(t *testing.T) {
	profile := &UserProfile{Nickname: "original"}

	invalid := strings.Repeat("a", ProfileMaxLength+1)
	err := profile.SetNickname(invalid)

	require.ErrorIs(t, err, errno.ErrInvalidProfileLength)
	require.Equal(t, "original", profile.Nickname)
}

func TestUserProfileSetEmptyValueUsesDefault(t *testing.T) {
	profile := &UserProfile{Avatar: "avatar.png"}

	err := profile.SetAvatar("")

	require.NoError(t, err)
	require.Empty(t, profile.Avatar)
}

func TestUserProfileSetGenderAcceptsFemale(t *testing.T) {
	gender := int8(2)
	profile := &UserProfile{}

	err := profile.SetGender(gender)

	require.NoError(t, err)
	require.Equal(t, gender, profile.Gender)
}

func TestUserRole(t *testing.T) {
	roleID := uint64(12)
	user := &User{RoleID: &roleID}

	require.Equal(t, "12", user.Role())
}
