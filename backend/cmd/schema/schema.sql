PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;

--------------------------------------------------
-- USERS (Employees)
--------------------------------------------------

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    is_active INTEGER NOT NULL DEFAULT 1,
    name TEXT NOT NULL,
    username TEXT NOT NULL UNIQUE,
    email TEXT UNIQUE,
    phone TEXT,
    address TEXT,
    id_type TEXT CHECK (
        id_type IN (
            'nid',
            'passport',
            'driving_license',
            'birth_certificate',
            'trade_license',
            'other'
        )
    ),
    id_number TEXT,
    image_url TEXT,
    password TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'staff' CHECK (
    role IN ('admin', 'manager', 'accounts', 'staff')),
    refresh_token TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--------------------------------------------------
-- CUSTOMERS
--------------------------------------------------

CREATE TABLE IF NOT EXISTS customers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    notes TEXT,
    name TEXT NOT NULL,
    phone TEXT NOT NULL,
    email TEXT,
    company TEXT,
    address TEXT,
    id_type TEXT CHECK (
        id_type IN (
            'nid',
            'passport',
            'driving_license',
            'birth_certificate',
            'trade_license',
            'other'
        )
    ),
    id_number TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_customers_name ON customers(name);

CREATE INDEX IF NOT EXISTS idx_customers_phone ON customers(phone);

--------------------------------------------------
-- STORAGE CONTRACTS
--------------------------------------------------

CREATE TABLE IF NOT EXISTS contracts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    notes TEXT,
    customer_id INTEGER NOT NULL,
    duration_type TEXT NOT NULL DEFAULT 'month' CHECK (
        duration_type IN ('day', 'week', 'month', 'year')
    ),
    duration INTEGER,
    price REAL,
    security_deposit REAL DEFAULT 0,
    estimated_value REAL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'active' CHECK (
        status IN ('active', 'completed', 'cancelled')
    ),
    ended_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (customer_id) REFERENCES customers(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_contract_customer ON contracts(customer_id);

--------------------------------------------------
-- STORED ITEMS
--------------------------------------------------

CREATE TABLE IF NOT EXISTS items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    notes TEXT,
    contract_id INTEGER,
    name TEXT NOT NULL,
    quantity_unit TEXT NOT NULL DEFAULT 'pcs' CHECK (
        quantity_unit IN ('pcs','g','kg','ton','ml','liter')
    ),
    quantity INTEGER NOT NULL DEFAULT 1,
    weight REAL, -- in KG
    width REAL,  -- Left to Right
    height REAL, -- Top to Bottom
    length REAL, -- Front to Back
    image_url TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (contract_id) REFERENCES contracts(id),    
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_items_contract ON items(contract_id);

--------------------------------------------------
-- SHIPMENTS
--------------------------------------------------

CREATE TABLE IF NOT EXISTS shipments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    shipment_type TEXT NOT NULL CHECK (
        shipment_type IN ('pickup', 'delivery')
    ),
    item_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,

    --------------------------------------------------
    -- DELIVERY FIELDS (nullable for pickup)
    --------------------------------------------------

    vehicle_number TEXT,
    driver_name TEXT,
    driver_phone TEXT,
    from_location TEXT,
    company_name TEXT,
    to_location TEXT,
    receiver_name TEXT,
    receiver_phone TEXT,
    customer_charge REAL DEFAULT 0,
    company_cost REAL DEFAULT 0,
    company_paid REAL DEFAULT 0,
    customer_paid REAL DEFAULT 0,

    --------------------------------------------------
    -- STATUS
    --------------------------------------------------

    status TEXT NOT NULL DEFAULT 'pending' CHECK (
        status IN (
            'pending',
            'in_transit',
            'completed',
            'cancelled'
        )
    ),

    --------------------------------------------------
    -- COMMON FIELDS
    --------------------------------------------------

    notes TEXT,
    shipped_at DATETIME,
    received_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (item_id) REFERENCES items(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

--------------------------------------------------
-- EXPENSES (DAILY BUSINESS COSTS)
--------------------------------------------------

CREATE TABLE IF NOT EXISTS expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    title TEXT NOT NULL,
    category TEXT NOT NULL CHECK (
        category IN (
            'rent',
            'salary',
            'bill',
            'other'
        )
    ),
    amount REAL NOT NULL,
    expense_date DATE NOT NULL,
    notes TEXT,
    image_url TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_expenses_date ON expenses(expense_date);

--------------------------------------------------
-- LOGS
--------------------------------------------------

CREATE TABLE IF NOT EXISTS logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,  -- who did the action
    action TEXT NOT NULL,  -- create, update, delete, approve, etc
    description TEXT NOT NULL,
    entity_type TEXT NOT NULL,  -- table name like 'items', 'contracts'
    entity_id INTEGER NOT NULL,  -- row id
    old_data TEXT,  -- JSON string
    new_data TEXT,  -- JSON string
    ip_address TEXT,
    user_agent TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id)
);









