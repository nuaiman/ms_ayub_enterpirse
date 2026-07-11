PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;

--------------------------------------------------
-- USERS (Employees)
-- Description: Stores system user accounts, employee information, and authentication details
--------------------------------------------------

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,          -- Unique user identifier
    name TEXT NOT NULL,                            -- Full name of the user
    username TEXT UNIQUE,                          -- Unique username for login
    password TEXT,                                 -- Hashed password
    email TEXT UNIQUE,                             -- User's email address
    phone TEXT,                                    -- User's phone number
    address TEXT,                                  -- User's physical address
    id_type TEXT CHECK (                           -- Type of identification document
        id_type IN ('nid', 'passport', 'driving_license', 'birth_certificate', 'trade_license', 'other')
    ),
    id_number TEXT,                                -- Identification document number
    image_url TEXT,                                -- URL to user's profile image
    role TEXT NOT NULL DEFAULT 'staff' CHECK (    -- User role for access control
        role IN ('admin', 'manager', 'accounts', 'staff')
    ),
    is_active INTEGER NOT NULL DEFAULT 1,          -- Whether user account is active (1=active, 0=inactive)
    monthly_salary REAL DEFAULT 0,                 -- Monthly salary amount for this user
    refresh_token TEXT,                            -- JWT refresh token for authentication
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- Record creation timestamp
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP   -- Record last update timestamp
);

--------------------------------------------------
-- ITEMS (Storage Items)
--------------------------------------------------

CREATE TABLE IF NOT EXISTS items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,          -- Unique item identifier
    user_id INTEGER NOT NULL,                      -- Employee who registered this item
    notes TEXT,                                    -- Additional notes about the item
    
    -- Customer Fields
    name TEXT NOT NULL,                            -- (Item's reference name) Customer's full name / company / organization name
    customer_phone TEXT NOT NULL,                  -- Customer's contact number
    customer_email TEXT,                           -- Customer's email address
    
    -- Category
    category TEXT,                                 -- Item category (free text with autocomplete)
    subcategory TEXT,                              -- Item subcategory (free text with autocomplete)
    
    -- Item Details
    quantity_unit TEXT NOT NULL DEFAULT 'pcs' CHECK (  -- Unit of measurement
        quantity_unit IN ('bag', 'bottle', 'box', 'can', 'carton', 'cup', 'dozen', 'gallon', 'pack', 'pair', 'pcs', 'roll', 'set', 'sheet', 'unit')
    ),
    quantity INTEGER NOT NULL DEFAULT 1,           -- Number of units
    
    -- Dimensions
    weight REAL,                                   -- Weight of the item
    weight_unit TEXT CHECK (weight_unit IN ('mg', 'g', 'oz', 'lb', 'kg', 'ton')) DEFAULT 'kg',  -- Weight unit
    width REAL,                                    -- Width measurement (Left to Right)
    height REAL,                                   -- Height measurement (Top to Bottom)
    length REAL,                                   -- Length measurement (Front to Back)
    dimension_unit TEXT CHECK (dimension_unit IN ('mm', 'cm', 'in', 'ft', 'm', 'yd', 'km')) DEFAULT 'in',  -- Dimension unit
    
    -- Storage Contract
    duration_type TEXT CHECK (duration_type IN ('day', 'week', 'month', 'year')),  -- Storage duration type
    duration INTEGER,                              -- Storage duration value
    start_date DATE,                               -- When storage started

    -- Revenue Details
    amount REAL NOT NULL,                          -- Total amount for this storage contract
    deposit REAL DEFAULT 0,                        -- Advanced deposit/security deposit
    customer_paid REAL DEFAULT 0,                  -- Amount customer has paid so far
    
    -- Status
    status TEXT DEFAULT 'active' CHECK (status IN ('active', 'complete')),  -- Storage status
    
    image_url TEXT,                                -- URL to item image
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- Record creation timestamp
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- Record last update timestamp

    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_items_user_id ON items(user_id);               -- For filtering items by user
CREATE INDEX IF NOT EXISTS idx_items_status ON items(status);                 -- For filtering by status
CREATE INDEX IF NOT EXISTS idx_items_customer_phone ON items(customer_phone); -- For searching by customer phone
CREATE INDEX IF NOT EXISTS idx_items_category ON items(category);             -- For category filtering
CREATE INDEX IF NOT EXISTS idx_items_subcategory ON items(subcategory);       -- For subcategory filtering
CREATE INDEX IF NOT EXISTS idx_items_start_date ON items(start_date);         -- For date-based queries

--------------------------------------------------
-- CHECKOUTS (Operations/Logistics)
-- Description: Tracks pickup and delivery operations for items
-- For pickup: only basic info + pickup details
-- For delivery: full details including location, transport, and financial
--------------------------------------------------

CREATE TABLE IF NOT EXISTS checkouts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,          -- Unique checkout identifier
    user_id INTEGER NOT NULL,                      -- User who created this checkout
    item_id INTEGER NOT NULL,                      -- Reference to the item being checked out
    quantity INTEGER NOT NULL DEFAULT 1,           -- Quantity in this checkout
    weight REAL DEFAULT 0;
    checkout_date DATETIME NOT NULL,               -- Date/time when checkout was created
    receiver_name TEXT,
    receiver_phone TEXT,
    
    -- Checkout Type
    type TEXT NOT NULL CHECK (type IN ('pickup', 'delivery')),  -- Type of checkout operation


    -- Delivery Details (only for delivery type)
    vehicle_number TEXT,
    driver_name TEXT,
    driver_phone TEXT,
    from_location TEXT,
    to_location TEXT,

    -- Financial (only for delivery)
    delivery_charge REAL DEFAULT 0,                -- Amount charged to customer for delivery
    delivery_cost REAL DEFAULT 0,                  -- Our cost for this delivery
    customer_paid REAL DEFAULT 0,                  -- Amount customer paid for this delivery
    
    notes TEXT,                                    -- Additional notes about the checkout
    image_url TEXT,                                -- URL to checkout photos/evidence
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- Record creation timestamp
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- Record last update timestamp
    
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_checkouts_user_id ON checkouts(user_id);       -- For filtering by user
CREATE INDEX IF NOT EXISTS idx_checkouts_item_id ON checkouts(item_id);       -- For finding checkouts for an item
CREATE INDEX IF NOT EXISTS idx_checkouts_type ON checkouts(type);             -- For filtering by type
CREATE INDEX IF NOT EXISTS idx_checkouts_status ON checkouts(status);         -- For filtering by status

--------------------------------------------------
-- EXPENSES (ALL cost tracking - SOURCE OF TRUTH)
--------------------------------------------------

CREATE TABLE IF NOT EXISTS expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,          -- Unique expense identifier
    user_id INTEGER NOT NULL,                      -- User who recorded this expense

    -- Salary Tracking
    is_salary BOOLEAN DEFAULT 0,                   -- Flag to identify salary payments
    salary_user_id INTEGER,                        -- User ID this salary is for (if is_salary is true)
    salary_month_year TEXT,                        -- Month and year of salary (Format: YYYY-MM)

    -- Expense Details
    title TEXT NOT NULL,                           -- Brief title/description of the expense
    category TEXT,                                 -- Expense category (free text with autocomplete)
    amount REAL NOT NULL,                          -- Total expense amount
    expense_date DATE NOT NULL,                    -- Date when expense was incurred
    notes TEXT,                                    -- Additional notes
    image_url TEXT,                                -- URL to bill/receipt image
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- Record creation timestamp
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- Record last update timestamp

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (salary_user_id) REFERENCES users(id),
    UNIQUE(salary_user_id, salary_month_year)      -- Ensure only one salary per user per month
);

CREATE INDEX IF NOT EXISTS idx_expenses_user_id ON expenses(user_id);               -- For filtering by user
CREATE INDEX IF NOT EXISTS idx_expenses_category ON expenses(category);             -- For category filtering
CREATE INDEX IF NOT EXISTS idx_expenses_date ON expenses(expense_date);             -- For date-based queries
CREATE INDEX IF NOT EXISTS idx_expenses_is_salary ON expenses(is_salary);           -- For salary filtering
CREATE INDEX IF NOT EXISTS idx_expenses_salary_user ON expenses(salary_user_id);    -- For finding salaries by user
CREATE INDEX IF NOT EXISTS idx_expenses_salary_month_year ON expenses(salary_month_year); -- For month-based salary queries

--------------------------------------------------
-- LOGS (Audit Trail)
--------------------------------------------------

CREATE TABLE IF NOT EXISTS logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,          -- Unique log identifier
    user_id INTEGER NOT NULL,                      -- User who performed the action
    action TEXT NOT NULL,                          -- Action type (CREATE, UPDATE, DELETE)
    description TEXT NOT NULL,                     -- Human-readable description of the action
    entity_type TEXT NOT NULL,                     -- Type of entity affected (user, item, checkout, expense)
    entity_id INTEGER NOT NULL,                    -- ID of the entity affected
    old_data TEXT,                                 -- Previous data before change (JSON format)
    new_data TEXT,                                 -- New data after change (JSON format)
    ip_address TEXT,                               -- IP address of the user who performed the action
    user_agent TEXT,                               -- Browser/device information
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- When the action occurred

    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_logs_user_id ON logs(user_id);               -- For filtering by user
CREATE INDEX IF NOT EXISTS idx_logs_entity_type ON logs(entity_type);       -- For filtering by entity type
CREATE INDEX IF NOT EXISTS idx_logs_entity_id ON logs(entity_id);           -- For finding logs for a specific entity
CREATE INDEX IF NOT EXISTS idx_logs_action ON logs(action);                -- For filtering by action type
CREATE INDEX IF NOT EXISTS idx_logs_created_at ON logs(created_at);         -- For date-based queries
