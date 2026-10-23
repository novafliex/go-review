package main

import (
	"context"
	"errors"
	"testing"
)

func TestExecTx(t *testing.T) {
	for _, fail := range []bool{false, true} {
		name := "commit"
		if fail {
			name = "rollback"
		}
		t.Run(name, func(t *testing.T) {
			store := &Store{db: &MockDB{}}
			sentinel := errors.New("business failed")
			var tx *MockTx
			err := store.execTx(context.Background(), func(current *MockTx) error {
				tx = current
				if fail {
					return sentinel
				}
				return nil
			})
			if fail && !errors.Is(err, sentinel) {
				t.Fatalf("error = %v", err)
			}
			if !fail && err != nil {
				t.Fatal(err)
			}
			if tx.committed != !fail || tx.rolledBack != fail {
				t.Fatalf("unexpected transaction state: %+v", tx)
			}
		})
	}
}
func TestExecTxPanicRollsBack(t *testing.T) {
	var tx *MockTx
	defer func() {
		if recover() != "boom" {
			t.Error("expected business panic")
		}
		if tx == nil || !tx.rolledBack || tx.committed {
			t.Error("transaction not rolled back")
		}
	}()
	store := &Store{db: &MockDB{}}
	_ = store.execTx(context.Background(), func(current *MockTx) error { tx = current; panic("boom") })
}
func TestTransferRejectsInvalidAmount(t *testing.T) {
	store := &Store{db: &MockDB{}}
	for _, amount := range []int{0, -1} {
		if store.TransferTx(context.Background(), TransferParams{Amount: amount}) == nil {
			t.Fatal("expected invalid amount error")
		}
	}
}
