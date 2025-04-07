package database

import (
	"context"
)

type NewerFunc[T any] func(db Database, tx *Transaction) *T

type Base[T any] struct {
	db      Database
	tx      *Transaction
	newFunc NewerFunc[T]
}

func NewBase[T any](db Database, tx *Transaction, fn NewerFunc[T]) *Base[T] {
	return &Base[T]{
		db:      db,
		tx:      tx,
		newFunc: fn,
	}
}

func (b *Base[T]) Begin() *T {
	return b.newFunc(b.db, b.tx)
}

func (b *Base[T]) BeginWithCtx(ctx context.Context) *T {
	return b.newFunc(b.db, b.db.TransactionWithContext(ctx))
}

func (b *Base[T]) HasTx() bool {
	return b.tx != nil
}

func (b *Base[T]) EnsureRollback() func() {
	if b.tx == nil {
		return func() {}
	}

	return b.tx.EnsureRollback
}

func (b *Base[T]) Commit() error {
	if b.tx == nil {
		return nil
	}

	return b.tx.Commit().Error
}
