package dao

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"mysql-practice/internal/model"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

// WithTx 返回一个使用事务连接的新 DAO（供 Service 层事务中使用）
func (d *UserDAO) WithTx(tx *gorm.DB) *UserDAO {
	return &UserDAO{db: tx}
}

// ---------------------------------------------------------------
// 基础 CRUD
// ---------------------------------------------------------------

// Create 插入单条记录
func (d *UserDAO) Create(ctx context.Context, user *model.User) error {
	return d.db.WithContext(ctx).Create(user).Error
}

// BatchCreate 批量插入
func (d *UserDAO) BatchCreate(ctx context.Context, users []model.User, batchSize int) error {
	return d.db.WithContext(ctx).CreateInBatches(&users, batchSize).Error
}

// GetByID 主键查询
func (d *UserDAO) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := d.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 唯一索引查询
func (d *UserDAO) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := d.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ListByAge 分页查询，返回列表 + 总数
func (d *UserDAO) ListByAge(ctx context.Context, minAge int, page, size int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := d.db.WithContext(ctx).Model(&model.User{}).Where("age >= ?", minAge)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Update 更新指定字段（用 map 解决零值问题）
func (d *UserDAO) Update(ctx context.Context, id uint, fields map[string]interface{}) error {
	result := d.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user %d not found", id)
	}
	return nil
}

// Delete 软删除
func (d *UserDAO) Delete(ctx context.Context, id uint) error {
	result := d.db.WithContext(ctx).Delete(&model.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user %d not found", id)
	}
	return nil
}

// ---------------------------------------------------------------
// Scopes：可复用的查询条件，链式调用
//
// 用法：db.Scopes(AgeRange(20, 30), Paginate(1, 10)).Find(&users)
//
// 企业中典型场景：
//   - 通用分页、排序
//   - 按状态过滤（上架/下架、启用/禁用）
//   - 时间范围（近7天、近30天）
//   - 权限过滤（只查自己部门的数据）
// ---------------------------------------------------------------

// AgeRange 年龄区间过滤
func AgeRange(min, max int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("age BETWEEN ? AND ?", min, max)
	}
}

// BalanceGte 余额大于等于（单位：分）
func BalanceGte(amount int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("balance >= ?", amount)
	}
}

// Paginate 通用分页，自带参数保护
func Paginate(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if size <= 0 {
			size = 10
		}
		if size > 100 {
			size = 100 // 防止一次拉太多数据打爆内存
		}
		return db.Offset((page - 1) * size).Limit(size)
	}
}

// OrderByCreated 按创建时间排序（desc=true 降序）
func OrderByCreated(desc bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if desc {
			return db.Order("created_at DESC")
		}
		return db.Order("created_at ASC")
	}
}

// ListWithScopes 使用 Scopes 组合查询（演示用法）
func (d *UserDAO) ListWithScopes(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) ([]model.User, error) {
	var users []model.User
	err := d.db.WithContext(ctx).Scopes(scopes...).Find(&users).Error
	return users, err
}
