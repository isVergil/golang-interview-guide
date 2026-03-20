package dao

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"mysql-practice/internal/model"
)

type OrderDAO struct {
	db *gorm.DB
}

func NewOrderDAO(db *gorm.DB) *OrderDAO {
	return &OrderDAO{db: db}
}

func (d *OrderDAO) WithTx(tx *gorm.DB) *OrderDAO {
	return &OrderDAO{db: tx}
}

// Create 创建订单
func (d *OrderDAO) Create(ctx context.Context, order *model.Order) error {
	return d.db.WithContext(ctx).Create(order).Error
}

// ---------------------------------------------------------------
// Preload vs Joins
//
// Preload：两次 SQL 查询，先查主表再查关联表，适合一对多
//   SELECT * FROM orders WHERE ...;
//   SELECT * FROM users WHERE id IN (1,2,3);
//
// Joins：单次 SQL 用 JOIN，适合一对一 / 需要用关联字段做 WHERE
//   SELECT orders.*, users.* FROM orders JOIN users ON ...;
// ---------------------------------------------------------------

// GetByOrderNo 按业务订单号查询（对外接口用这个）
func (d *OrderDAO) GetByOrderNo(ctx context.Context, orderNo int64) (*model.Order, error) {
	var order model.Order
	err := d.db.WithContext(ctx).Preload("User").Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByID 按自增主键查询（内部用）
func (d *OrderDAO) GetByID(ctx context.Context, id uint) (*model.Order, error) {
	var order model.Order
	err := d.db.WithContext(ctx).Preload("User").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// ListByUserID Preload 方式查用户的所有订单
func (d *OrderDAO) ListByUserID(ctx context.Context, userID uint) ([]model.Order, error) {
	var orders []model.Order
	err := d.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("User").
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// ListByUserName Joins 方式：用关联表字段做过滤条件
func (d *OrderDAO) ListByUserName(ctx context.Context, userName string) ([]model.Order, error) {
	var orders []model.Order
	err := d.db.WithContext(ctx).
		Joins("User"). // GORM v2 Joins preload
		Where("User.name = ?", userName).
		Find(&orders).Error
	return orders, err
}

// ListPaidWithUser 组合 Scopes + Preload
func (d *OrderDAO) ListPaidWithUser(ctx context.Context, page, size int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := d.db.WithContext(ctx).Model(&model.Order{}).Where("status = ?", 1)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Scopes(Paginate(page, size)).
		Preload("User").
		Order("created_at DESC").
		Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// UpdateStatus 更新订单状态
func (d *OrderDAO) UpdateStatus(ctx context.Context, id uint, status int8) error {
	result := d.db.WithContext(ctx).Model(&model.Order{}).
		Where("id = ?", id).
		Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("order %d not found", id)
	}
	return nil
}
