package main

import (
	"context"
	"errors"
	"fmt"
)

// 模拟事务对象
type MockTx struct {
	commited   bool
	rolledBack bool
}

func (tx *MockTx) Commit() error {
	tx.commited = true
	fmt.Println("[事务日志]:成功提交!")
	return nil
}

func (tx *MockTx) Rollback() error {
	tx.rolledBack = true
	fmt.Println("[事务日志]:发生错误,触发回滚!")
	return nil
}

// 模拟数据库连接
type MockDB struct{}

func (db *MockDB) BeginTx(ctx context.Context) (*MockTx, error) {
	fmt.Println("[事务日志]:开启新事物!")
	return &MockTx{}, nil
}

// store 结构体
type Store struct {
	db *MockDB
}

// 核心事务模板 execTx
func (store *Store) execTx(ctx context.Context, fn func(tx *MockTx) error) error {
	tx, err := store.db.BeginTx(ctx)
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// 业务逻辑模拟 transferTx
type TransferParams struct {
	Amount int
	IsFail bool
}

func (store *Store) TransferTx(ctx context.Context, arg TransferParams) error {
	return store.execTx(ctx, func(tx *MockTx) error {
		fmt.Println("--步骤--")

		if arg.IsFail {
			return errors.New("模拟出错")
		}

		fmt.Println("--步骤--")
		return nil
	})
}

func main() {
	store := &Store{db: &MockDB{}}
	ctx := context.Background()

	fmt.Println("--- 练习测试 1：正常流程 ---")
	_ = store.TransferTx(ctx, TransferParams{Amount: 100, IsFail: false})

	fmt.Println("\n--- 练习测试 2：异常回滚流程 ---")
	_ = store.TransferTx(ctx, TransferParams{Amount: 100, IsFail: true})
}
