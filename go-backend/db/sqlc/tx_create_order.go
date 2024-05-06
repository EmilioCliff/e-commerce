package db

import "context"

type TxCreateOrderParams struct {
}

type TxCreateOrderResult struct {
}

func (store *Store) CreateOrderTx(ctx context.Context, arg TxCreateOrderParams) (TxCreateOrderResult, error) {
	var result TxCreateOrderResult
	err := store.execTx(ctx, func(*Queries) error {
		return nil
	})
	return result, err
}
