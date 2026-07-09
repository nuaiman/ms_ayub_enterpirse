package models

import (
	"database/sql"
)

type Models struct {
	User     UserModel
	Item     ItemModel
	Checkout CheckoutModel
	Expense  ExpenseModel
	Log      LogModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		User:     UserModel{DB: db},
		Item:     ItemModel{DB: db},
		Checkout: CheckoutModel{DB: db},
		Expense:  ExpenseModel{DB: db},
		Log:      LogModel{DB: db},
	}
}
