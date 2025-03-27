package model

type UsersModel struct {
	// gorm.Model // 引入模板结构体  ID, CreatedAt, UpdatedAt, DeletedAt
	ID           int64  `json:"id" gorm:"primary_key;column:id"`
	UserName     int64  `json:"username" gorm:"column:username"`
	Email        string `json:"email" gorm:"column:email"`
	PasswordHash string `json:"password_hash" gorm:"column:password_hash"`
	//Status        int8    `json:"status" gorm:"column:status"`
	//EntBalance    float64 `json:"ent_balance" gorm:"column:ent_balance"`
	//ArtBalance    float64 `json:"art_balance" gorm:"column:art_balance"`
	//UserType      int8    `json:"user_type" gorm:"column:user_type"`
	//LoginIP       string  `json:"login_ip" gorm:"column:login_ip"`
	CreateAt  string              `json:"create_time" gorm:"column:create_at"`
	UpdateAt  string              `json:"update_time" gorm:"column:update_at"`
	BaseModel `json:"-" gorm:"-"` // 继承基础模型
}
