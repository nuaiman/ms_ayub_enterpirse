# MS Ayub Enterprise

A full-stack inventory and business management system built with Go (Chi router) and Vue 3 (TypeScript).

## Features

### User Management

- Role-based access control (Admin, Manager, Staff, Accounts)
- JWT authentication with refresh tokens
- Profile management with image upload
- Password reset and account activation/deactivation

### Customers

- CRUD operations with search, sort, and export
- ID type validation (NID, Passport, Driving License, etc.)
- Contact details management

### Contracts

- Create and manage contracts linked to customers
- Duration types: Day, Week, Month, Year
- Pricing, security deposit, and estimated value tracking
- Status management: Active, Completed, Cancelled
- Auto timestamp on status changes

### Items

- Inventory tracking with quantity and unit management
- Dimensions (L×W×H) and weight tracking
- Contract assignment for tracking item ownership
- Image upload support

### Shipments

- Pickup and delivery tracking
- Item quantity deduction on shipment creation
- Quantity restoration on cancellation/deletion
- Delivery details: Vehicle, Driver, Route, Receiver
- Financial tracking: Customer charge, Company cost

### Expenses

- Category-based expense tracking (Rent, Salary, Bill, Other)
- Receipt image upload
- Month/year filtering
- Total amount summaries

### Audit Logs

- Automatic logging of all create, update, and delete operations
- Entity type filtering
- Old/new data comparison with JSON diff view
- IP address and user agent tracking
- User-friendly display with username lookup

### Backup

- One-click database and file backup download
- Includes SQLite database and all uploaded files

## Tech Stack

### Backend

- **Language:** Go 1.21+
- **Router:** Chi v5
- **Database:** SQLite
- **Authentication:** JWT (access + refresh tokens)
- **File Storage:** Local filesystem

### Frontend

- **Framework:** Vue 3 (Composition API)
- **Language:** TypeScript
- **State Management:** Pinia
- **HTTP Client:** Axios
- **Styling:** Tailwind CSS v4
- **Build Tool:** Vite

## Getting Started

### Prerequisites

- Go 1.25+
- Node.js 18+
- npm or yarn

### Backend Setup

1. Navigate to the backend directory:
   cd backend
   go mod tidy
   go run ./cmd/

2. Navigate to the frontend directory:
   npm i
   npm run dev

### Role Permissions

Feature Admin Manager Staff Accounts
Customers ✅ ✅ ❌ ❌
Contracts ✅ ✅ ❌ ❌
Items ✅ ✅ ✅ ❌
Shipments ✅ ✅ ❌ ❌
Expenses ✅ ❌ ❌ ✅
Users ✅ ✅ ❌ ❌
Audit Logs ✅ ❌ ❌ ❌
Backup ✅ ❌ ❌ ❌
