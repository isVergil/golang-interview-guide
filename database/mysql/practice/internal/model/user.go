package model

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Name    string `gorm:"size:100;not null;comment:用户名"`
	Age     int    `gorm:"default:0;comment:年龄"`
	Email   string `gorm:"size:128;uniqueIndex;comment:邮箱"`
	Balance int64  `gorm:"default:0;comment:余额(分)"`
}

func (User) TableName() string {
	return "users"
}

// ---------------------------------------------------------------
// Hooks：GORM 在特定生命周期自动调用
// 适合放：数据清洗、自动填充、校验、审计日志
// 不适合放：业务逻辑（应放在 Service 层）
// ---------------------------------------------------------------

// BeforeCreate 创建前自动执行
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Name = strings.TrimSpace(u.Name)
	return nil
}

// BeforeUpdate 更新前自动执行
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Name = strings.TrimSpace(u.Name)
	return nil
}

// Order 订单模型（双 ID 模式：自增主键 + 雪花业务号）
type Order struct {
	ID          uint           `gorm:"primaryKey;comment:自增主键"`
	OrderNo     int64          `gorm:"uniqueIndex;not null;comment:订单号(雪花ID)"`
	UserID      uint           `gorm:"index;not null;comment:用户ID"`
	User        User           `gorm:"foreignKey:UserID"` // belongs to
	ProductName string         `gorm:"size:200;not null;comment:商品名"`
	Amount      int64          `gorm:"not null;comment:金额(分)"`
	Status      int8           `gorm:"default:0;comment:状态 0-待支付 1-已支付 2-已取消"`
	CreatedAt   time.Time      `gorm:"comment:创建时间"`
	UpdatedAt   time.Time      `gorm:"comment:更新时间"`
	DeletedAt   gorm.DeletedAt `gorm:"index;comment:删除时间"`
}

func (Order) TableName() string {
	return "orders"
}

// Product 商品模型（演示 Upsert + 乐观锁）
type Product struct {
	ID        uint           `gorm:"primaryKey"`
	SKU       string         `gorm:"size:64;uniqueIndex;not null;comment:商品编码"`
	Name      string         `gorm:"size:200;not null;comment:商品名"`
	Price     int64          `gorm:"not null;comment:价格(分)"`
	Stock     int            `gorm:"default:0;comment:库存"`
	Version   int            `gorm:"default:0;comment:乐观锁版本号"`
	CreatedAt time.Time      `gorm:"comment:创建时间"`
	UpdatedAt time.Time      `gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间"`
}

func (Product) TableName() string {
	return "products"
}
