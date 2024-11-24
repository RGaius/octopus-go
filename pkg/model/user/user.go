package user

type User struct {
	ID         int    `json:"id" xorm:"id"`
	Username   string `json:"username" xorm:"username"`
	Password   string `json:"password" xorm:"password"`
	CreateDate string `json:"create_date" xorm:"create_date"`
	UpdateDate string `json:"update_date" xorm:"update_date"`
}
