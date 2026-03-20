package dao

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"mysql-practice/internal/model"
)

type ProductDAO struct {
	db *gorm.DB
}

func NewProductDAO(db *gorm.DB) *ProductDAO {
	return &ProductDAO{db: db}
}

// ---------------------------------------------------------------
// Upsert：INSERT ... ON DUPLICATE KEY UPDATE
//
// 场景：幂等写入，比如同步外部数据、消息消费去重
// 按 SKU 唯一键冲突时更新价格和库存，不重复插入
// ---------------------------------------------------------------

// Upsert 按 SKU 幂等写入，冲突时更新 price 和 stock
func (d *ProductDAO) Upsert(ctx context.Context, product *model.Product) error {
	return d.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "sku"}},                         // 冲突判断列
			DoUpdates: clause.AssignmentColumns([]string{"name", "price", "stock"}), // 冲突时更新的列
		}).
		Create(product).Error
}

// BatchUpsert 批量幂等写入
func (d *ProductDAO) BatchUpsert(ctx context.Context, products []model.Product, batchSize int) error {
	return d.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "sku"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "price", "stock"}),
		}).
		CreateInBatches(&products, batchSize).Error
}

// GetBySKU 按 SKU 查询
func (d *ProductDAO) GetBySKU(ctx context.Context, sku string) (*model.Product, error) {
	var product model.Product
	err := d.db.WithContext(ctx).Where("sku = ?", sku).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// ---------------------------------------------------------------
// 乐观锁：基于 version 字段的并发安全更新
//
// 原理：
//   UPDATE products SET stock = ?, version = version + 1
//   WHERE id = ? AND version = ?
//
// 如果 RowsAffected == 0，说明别人已经改过了，当前操作失败
//
// 对比悲观锁（FOR UPDATE）：
//   - 悲观锁：加行锁，阻塞其他事务，适合冲突率高的场景
//   - 乐观锁：不加锁，用版本号检测冲突，适合冲突率低的场景
// ---------------------------------------------------------------

// DeductStockOptimistic 乐观锁扣库存
// 返回 error == nil 表示扣减成功
func (d *ProductDAO) DeductStockOptimistic(ctx context.Context, id uint, quantity int, currentVersion int) error {
	result := d.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("id = ? AND version = ? AND stock >= ?", id, currentVersion, quantity).
		Updates(map[string]interface{}{
			"stock":   gorm.Expr("stock - ?", quantity),
			"version": gorm.Expr("version + 1"),
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("optimistic lock conflict or insufficient stock (id=%d, version=%d)", id, currentVersion)
	}
	return nil
}
