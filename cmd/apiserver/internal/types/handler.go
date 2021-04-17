package types

type CreateUserReq struct {
	Name     string `json:"name" binding:"required"`
	Mobile   string `json:"mobile" binding:"required,checkMobile"`
	Password string `json:"password" binding:"required,gte=6"`
}
