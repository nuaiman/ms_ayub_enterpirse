package main

import (
	"backend/internal/app"
	"backend/internal/bootstrap"
	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/handlers"
	"backend/internal/models"
	"backend/internal/router"
	"log"
)

func main() {
	cfg := config.MustLoadConfig()

	dbPool := db.InitDB(cfg.DBPath, cfg.DBName, cfg.SchemaPath)
	defer db.CloseDB(dbPool)

	a := app.Application{
		Config: cfg,
		DB:     dbPool,
		Models: models.NewModel(dbPool),
	}

	bootstrap.BootstrapAdmin(&a)

	h := handlers.New(&a)

	r := router.RegisterRouter(&a, h)

	log.Fatal(a.Run(r))
}

// go run ./cmd/app/ --config=config/config.env

// next lets move

// CREATE TABLE IF NOT EXISTS customers (
//     id INTEGER PRIMARY KEY AUTOINCREMENT,
//     user_id INTEGER,
//     is_approved INTEGER NOT NULL DEFAULT 0,
//     notes TEXT,
//     name TEXT NOT NULL,
//     phone TEXT NOT NULL,
//     email TEXT,
//     company TEXT,
//     address TEXT,
//     id_type TEXT CHECK (
//         id_type IN (
//             'nid',
//             'passport',
//             'driving_license',
//             'birth_certificate',
//             'trade_license',
//             'other'
//         )
//     ),
//     id_number TEXT,
//     created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
//     updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
//     FOREIGN KEY (user_id) REFERENCES users(id)
// );

//  we will need create getByID GetAll Update
//  update is manager or admin only
//  is_approved is admin only
