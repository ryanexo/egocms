package model

import (
	"net/mail"
	"regexp"
	"strconv"
	"time"

	"cms/internal/app/user/errno"
	"cms/internal/infra/store/modeltype"
)

const (
	UsernameMinLength = 4
	UsernameMaxLength = 255
	EmailMaxLength    = 255
	ProfileMaxLength  = 255
)

var usernamePattern = regexp.MustCompile(`^[A-Za-z0-9]+$`)

type User struct {
	modeltype.Base
	Username   string
	Password   string
	Email      string
	VerifiedAt *time.Time
	IP         *string
	Status     int8
	RoleID     *uint64
	RoleName   string
	Profile    *UserProfile `gorm:"foreignKey:UserID;references:ID"`
}

func (User) TableName() string { return "user" }

func (user *User) SetUsername(username string) error {
	if username == "" {
		return errno.ErrUsernameRequired
	}
	if len(username) < UsernameMinLength || len(username) > UsernameMaxLength || !usernamePattern.MatchString(username) {
		return errno.ErrInvalidUsername
	}
	user.Username = username
	return nil
}

func (user *User) SetEmail(email string) error {
	if email == "" {
		return errno.ErrEmailRequired
	}
	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email || len(email) > EmailMaxLength {
		return errno.ErrInvalidEmail
	}
	user.Email = email
	return nil
}

func (user *User) SetPasswordHash(passwordHash string) error {
	if passwordHash == "" {
		return errno.ErrPasswordRequired
	}
	user.Password = passwordHash
	return nil
}

func (user *User) UserID() uint64 { return user.ID }

func (user *User) Role() string {
	if user.RoleID == nil {
		return ""
	}
	return strconv.FormatUint(*user.RoleID, 10)
}

type UserProfile struct {
	modeltype.Base
	Avatar      string
	UserID      uint64
	Nickname    string
	Gender      int8
	Description string
	Country     string
	Province    string
	City        string
}

func (profile *UserProfile) SetAvatar(value string) error { return setString(&profile.Avatar, value) }

func (profile *UserProfile) SetNickname(value string) error {
	return setString(&profile.Nickname, value)
}

func (profile *UserProfile) SetDescription(value string) error {
	return setString(&profile.Description, value)
}

func (profile *UserProfile) SetCountry(value string) error {
	return setString(&profile.Country, value)
}

func (profile *UserProfile) SetProvince(value string) error {
	return setString(&profile.Province, value)
}

func (profile *UserProfile) SetCity(value string) error { return setString(&profile.City, value) }

func (profile *UserProfile) SetGender(value int8) error {
	if value < 0 || value > 2 {
		return errno.ErrInvalidGender
	}
	profile.Gender = value
	return nil
}

func setString(target *string, value string) error {
	if len(value) > ProfileMaxLength {
		return errno.ErrInvalidProfileLength
	}
	*target = value
	return nil
}
