PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;

--------------------------------------------------
-- USERS (Employees)
--------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    username TEXT UNIQUE,
    password TEXT,
    email TEXT UNIQUE,
    phone TEXT,
    address TEXT,
    id_type TEXT CHECK (
        id_type IN ('nid', 'passport', 'driving_license', 'birth_certificate', 'trade_license', 'other')
    ),
    id_number TEXT,
    image_url TEXT,
    role TEXT NOT NULL DEFAULT 'staff' CHECK (
        role IN ('admin', 'manager', 'accounts', 'staff')
    ),
    is_active INTEGER NOT NULL DEFAULT 1,
    monthly_salary REAL DEFAULT 0,
    refresh_token TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--------------------------------------------------
-- ITEMS (Storage Items)
--------------------------------------------------
CREATE TABLE IF NOT EXISTS items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    notes TEXT,
    
    -- Product & Storage Info
    product_name TEXT NOT NULL,
    storage_name TEXT,
    account_name TEXT,
    lot_number TEXT,
    
    -- Customer Fields
    customer_phone TEXT,                  -- Made optional
    customer_email TEXT,
    
    -- Category
    category TEXT,
    subcategory TEXT,
    
    -- Item Details
    quantity_unit TEXT NOT NULL DEFAULT 'pcs' CHECK (
        quantity_unit IN ('bag', 'bottle', 'box', 'can', 'carton', 'cup', 'dozen', 'gallon', 'pack', 'pair', 'pcs', 'roll', 'set', 'sheet', 'unit')
    ),
    quantity INTEGER NOT NULL DEFAULT 1,
    
    -- Revenue Details
    amount REAL DEFAULT 0,
    deposit REAL DEFAULT 0,
    customer_paid REAL DEFAULT 0,
    
    image_url TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_items_user_id ON items(user_id);
CREATE INDEX IF NOT EXISTS idx_items_customer_phone ON items(customer_phone);
CREATE INDEX IF NOT EXISTS idx_items_category ON items(category);
CREATE INDEX IF NOT EXISTS idx_items_subcategory ON items(subcategory);
CREATE INDEX IF NOT EXISTS idx_items_product_name ON items(product_name);
CREATE INDEX IF NOT EXISTS idx_items_lot_number ON items(lot_number);

--------------------------------------------------
-- CHECKOUTS (Operations/Logistics)
-- No status column, weight added
--------------------------------------------------
CREATE TABLE IF NOT EXISTS checkouts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    item_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    weight REAL DEFAULT 0,             -- ✅ fixed: comma separated, no rogue semicolon
    checkout_date DATETIME NOT NULL,
    receiver_name TEXT,
    receiver_phone TEXT,
    type TEXT NOT NULL CHECK (type IN ('pickup', 'delivery')),
    vehicle_number TEXT,
    driver_name TEXT,
    driver_phone TEXT,
    from_location TEXT,
    to_location TEXT,
    delivery_charge REAL DEFAULT 0,
    delivery_cost REAL DEFAULT 0,
    customer_paid REAL DEFAULT 0,
    notes TEXT,
    image_url TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE RESTRICT
);

-- Indexes: removed idx_checkouts_status since status column was dropped
CREATE INDEX IF NOT EXISTS idx_checkouts_user_id ON checkouts(user_id);
CREATE INDEX IF NOT EXISTS idx_checkouts_item_id ON checkouts(item_id);
CREATE INDEX IF NOT EXISTS idx_checkouts_type ON checkouts(type);

--------------------------------------------------
-- EXPENSES (ALL cost tracking)
--------------------------------------------------
CREATE TABLE IF NOT EXISTS expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    is_salary BOOLEAN DEFAULT 0,
    salary_user_id INTEGER,
    salary_month_year TEXT,
    title TEXT NOT NULL,
    category TEXT,
    amount REAL NOT NULL,
    expense_date DATE NOT NULL,
    notes TEXT,
    image_url TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (salary_user_id) REFERENCES users(id),
    UNIQUE(salary_user_id, salary_month_year)
);

CREATE INDEX IF NOT EXISTS idx_expenses_user_id ON expenses(user_id);
CREATE INDEX IF NOT EXISTS idx_expenses_category ON expenses(category);
CREATE INDEX IF NOT EXISTS idx_expenses_date ON expenses(expense_date);
CREATE INDEX IF NOT EXISTS idx_expenses_is_salary ON expenses(is_salary);
CREATE INDEX IF NOT EXISTS idx_expenses_salary_user ON expenses(salary_user_id);
CREATE INDEX IF NOT EXISTS idx_expenses_salary_month_year ON expenses(salary_month_year);

--------------------------------------------------
-- LOGS (Audit Trail)
--------------------------------------------------
CREATE TABLE IF NOT EXISTS logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    action TEXT NOT NULL,
    description TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id INTEGER NOT NULL,
    old_data TEXT,
    new_data TEXT,
    ip_address TEXT,
    user_agent TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_logs_user_id ON logs(user_id);
CREATE INDEX IF NOT EXISTS idx_logs_entity_type ON logs(entity_type);
CREATE INDEX IF NOT EXISTS idx_logs_entity_id ON logs(entity_id);
CREATE INDEX IF NOT EXISTS idx_logs_action ON logs(action);
CREATE INDEX IF NOT EXISTS idx_logs_created_at ON logs(created_at);