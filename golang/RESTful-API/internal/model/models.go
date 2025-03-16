package model // 定义包名为 model，用于数据库模型

import (
	"gorm.io/gorm" // 引入 GORM 包，用于数据库操作
	"time"         // 引入 time 包，用于处理时间相关的字段
)

// CrudList gorm结构体
type CrudList struct {
	gorm.Model // 引入 GORM 的基础模型，自动包含 ID, CreatedAt, UpdatedAt 和 DeletedAt 字段
	//ID        uint `gorm:"primarykey"`  // 如果不使用 gorm.Model 可以手动定义 ID 字段
	//CreatedAt time.Time  // 自动包含 CreatedAt 字段
	//UpdatedAt time.Time  // 自动包含 UpdatedAt 字段
	//DeletedAt DeletedAt `gorm:"index"` // 自动包含 DeletedAt 字段，用于软删除

	Name     string `gorm:"type: varchar(50); not null" binding:"required" json:"name"` // 定义 Name 字段，类型为 varchar(50)，不可为空，绑定为必填字段，返回时为 JSON 中的 "name"
	Level    string `gorm:"type: varchar(20);" json:"level"`                            // 定义 Level 字段，类型为 varchar(20)，返回时为 JSON 中的 "level"
	Email    string `gorm:"type: varchar(50);" json:"email"`                            // 定义 Email 字段，类型为 varchar(50)，返回时为 JSON 中的 "email"
	Phone    string `gorm:"type: varchar(20);" json:"phone"`                            // 定义 Phone 字段，类型为 varchar(20)，返回时为 JSON 中的 "phone"
	Birthday string `gorm:"type: varchar(20);" json:"birthday"`                         // 定义 Birthday 字段，类型为 varchar(20)，返回时为 JSON 中的 "birthday"
	Address  string `gorm:"type: varchar(200);" json:"address"`                         // 定义 Address 字段，类型为 varchar(200)，返回时为 JSON 中的 "address"
	// 变量名要大写，才public，可以被gorm获取解析
}

type UserList struct {
	Name        string    `gorm:"type: varchar(50); not null; primarykey" binding:"required" json:"name"` // 定义 Name 字段，类型为 varchar(50)，不可为空，作为主键，绑定为必填字段，返回时为 JSON 中的 "name"
	AuthGroup   string    `gorm:"type: varchar(20);" json:"auth_level"`                                   // 定义 AuthGroup 字段，类型为 varchar(20)，返回时为 JSON 中的 "auth_level"
	Email       string    `gorm:"type: varchar(50); not null;" json:"email" binding:"required"`           // 定义 Email 字段，类型为 varchar(50)，不可为空，绑定为必填字段，返回时为 JSON 中的 "email"
	Phone       string    `gorm:"type: varchar(20);" json:"phone"`                                        // 定义 Phone 字段，类型为 varchar(20)，返回时为 JSON 中的 "phone"
	Password    string    `gorm:"type: varchar(100); not null;" json:"password" binding:"required"`       // 定义 Password 字段，类型为 varchar(100)，不可为空，绑定为必填字段，返回时为 JSON 中的 "password"
	PasswordTry uint      `json:"password_try"`                                                           // 定义 PasswordTry 字段，类型为 uint，返回时为 JSON 中的 "password_try"
	LockedUntil time.Time `json:"locked_until"`                                                           // 定义 LockedUntil 字段，类型为 time.Time，返回时为 JSON 中的 "locked_until"
	Status      string    `json:"status"`                                                                 // 定义 Status 字段，类型为 string，返回时为 JSON 中的 "status"
}
