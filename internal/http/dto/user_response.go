package dto

type UserInfo struct {
    ID       int64  `json:"id"`
    Username string `json:"username"`
    RoleID   int    `json:"roleId"`
}

type UserLoginResponse struct {
    Token string `json:"token"`
}
