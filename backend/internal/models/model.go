package models

import (
	"database/sql"
)

type Models struct {
	User     UserModel
	Customer CustomerModel
	Contract ContractModel
	Item     ItemModel
	Shipment ShipmentModel
	Expense  ExpenseModel
	Log      LogModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		User:     UserModel{DB: db},
		Customer: CustomerModel{DB: db},
		Contract: ContractModel{DB: db},
		Item:     ItemModel{DB: db},
		Shipment: ShipmentModel{DB: db},
		Expense:  ExpenseModel{DB: db},
		Log:      LogModel{DB: db},
	}
}
