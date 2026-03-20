package service

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"mysql-practice/internal/dao"
	"mysql-practice/internal/model"
)

type UserService struct {
	db      *gorm.DB
	userDAO *dao.UserDAO
}

func NewUserService(db *gorm.DB, userDAO *dao.UserDAO) *UserService {
	return &UserService{db: db, userDAO: userDAO}
}

// CreateUser 创建用户（参数校验 + 调用 DAO）
func (s *UserService) CreateUser(ctx context.Context, name string, age int, email string) (*model.User, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	user := &model.User{
		Name:  name,
		Age:   age,
		Email: email,
	}

	if err := s.userDAO.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

// Transfer 转账：在事务中完成 from 扣款 + to 加款
func (s *UserService) Transfer(ctx context.Context, fromID, toID uint, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 事务内创建临时 DAO，使用 tx 而非原始 db
		txDAO := s.userDAO.WithTx(tx)

		// 查询转出方（SELECT ... FOR UPDATE 加行锁防止并发）
		var from model.User
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&from, fromID).Error; err != nil {
			return fmt.Errorf("query from user: %w", err)
		}

		// 余额检查
		if from.Balance < amount {
			return fmt.Errorf("insufficient balance: have %d, need %d", from.Balance, amount)
		}

		// 查询转入方
		var to model.User
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&to, toID).Error; err != nil {
			return fmt.Errorf("query to user: %w", err)
		}

		// 扣款
		if err := txDAO.Update(ctx, fromID, map[string]interface{}{
			"balance": gorm.Expr("balance - ?", amount),
		}); err != nil {
			return fmt.Errorf("deduct from user: %w", err)
		}

		// 加款
		if err := txDAO.Update(ctx, toID, map[string]interface{}{
			"balance": gorm.Expr("balance + ?", amount),
		}); err != nil {
			return fmt.Errorf("add to user: %w", err)
		}

		return nil // 返回 nil 自动 Commit，返回 error 自动 Rollback
	})
}
