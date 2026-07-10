# AkadiaPay

> School Financial Platform for configurable student billing, payment processing, and financial ledger management.

AkadiaPay is a backend MVP for a **School Financial Platform** designed to help schools manage student financial obligations, billing configuration, payment processing, allocation, and financial ledger records in a structured and auditable way.

This project started as a final project for a Golang backend bootcamp, but the product direction is larger than a bootcamp assignment. AkadiaPay is part of the broader **Akadia** vision: a modular school financial ecosystem that can grow into a commercial SaaS product.

---

## Project Status

This public repository is now **frozen**.

The bootcamp phase has been completed, and this repository represents the public MVP snapshot used for learning, demonstration, and portfolio purposes.

Future development of Akadia will continue in a **private repository** because Akadia is intended to become a real product that can be developed, monetized, and maintained professionally by the owner and company behind it.

This repository will not represent the full commercial roadmap.

---

## Overview

Many schools still manage student payments manually using spreadsheets, messaging apps, and disconnected payment records. This makes billing difficult to track, payment rules hard to customize, and reconciliation time-consuming.

AkadiaPay addresses this problem by providing a backend system where schools can:

- Configure payment rules.
- Create payment products.
- Generate student obligations.
- Process payment orders.
- Allocate payments to one or more obligations.
- Track outstanding balances.
- Record financial ledger entries.
- Support multi-tenant school operations.

The central idea is simple:

```text
Login
  ↓
Payment Policy
  ↓
Payment Product
  ↓
Student Obligation
  ↓
Payment Order
  ↓
Payment Allocation
  ↓
Financial Ledger
```

---

## Core Concept

AkadiaPay is designed around the idea that **Student Obligation is the center of the payment domain**.

Payment is not the main object.

Payment is an event that reduces a student's financial obligation.

```text
Student Obligation = What the student owes
Payment Order      = Payment transaction
Allocation         = How payment reduces obligations
Ledger Entry       = Accounting record of the transaction
```

This design makes the system more flexible, auditable, and easier to extend into reporting, settlement, and accounting modules.

---

## Key Principles

### 1. Multi-Tenant

AkadiaPay supports multiple schools in one application infrastructure.

Each school has isolated:

- Users
- Students
- Payment policies
- Payment products
- Student obligations
- Payment orders
- Ledger entries

Tenant isolation is a core security and business requirement.

### 2. Rule-Based Billing

Schools can configure payment behavior using payment policies.

Examples:

- Full payment only
- Partial payment allowed
- Minimum payment amount
- Minimum payment percentage
- Installment-style payment
- Auto-close obligation after full payment

### 3. Student Obligation Centric

The system focuses on what the student owes, not only on payment transactions.

This makes it easier to track:

- Original amount
- Paid amount
- Outstanding amount
- Partial payment
- Closed obligation
- Cancelled obligation

### 4. Ledger-Based Financial Tracking

Every completed payment can be represented as debit and credit ledger entries.

Example:

```text
Debit  Cash / Bank          Rp500.000
Credit Tuition Revenue      Rp500.000
```

This provides a foundation for reconciliation, audit trail, and future accounting integration.

### 5. Modular Architecture

AkadiaPay is built as a monolith-first backend, but with modular boundaries that make future extraction possible.

Major modules:

- Authentication & Identity
- Master Payment Configuration
- Student Billing
- Payment Processing
- Financial Ledger
- Reporting foundation

---

## MVP Scope

The MVP covers the following bootcamp and product requirements:

- JWT Authentication
- Role-Based Access Control
- Multi-tenant login context
- CRUD for payment policies
- CRUD for payment products
- Student obligation management
- Payment order creation
- Payment allocation
- Financial ledger posting
- Pagination and filtering
- Seed data for demo scenarios
- Bruno REST API collection
- Swagger/OpenAPI documentation
- PlantUML business-process documentation
- Docker-ready backend setup
- Clean Architecture-inspired layering

---

## Main Business Flow

### 1. Authentication

Users log in using email and password.  
The system identifies the active tenant, role, and user context.

Supported roles include:

- Super Admin
- School Admin
- Treasurer / Finance Officer
- Parent
- Student

### 2. Payment Policy

A school defines **how a bill can be paid**.

Examples:

- Full payment
- Partial payment
- Minimum amount
- Minimum percentage
- Installment-like behavior

### 3. Payment Product

A school defines **what will be billed**.

Examples:

- SPP
- Registration fee
- Building fee
- Study tour
- Book package
- Uniform fee

Each payment product is connected to a payment policy.

### 4. Student Obligation

A payment product becomes a real bill when assigned to a student.

Example:

```text
Student: Kevin Wijaya
Product: SPP August 2026
Original Amount: Rp500.000
Outstanding Amount: Rp500.000
Status: Pending
```

### 5. Payment Order

A parent or finance officer creates a payment order.

The payment order represents the payment transaction.

### 6. Payment Allocation

The payment is allocated to one or more student obligations.

This process updates:

- Paid amount
- Outstanding amount
- Obligation status
- Payment order status

### 7. Financial Ledger

After payment is completed, the system records ledger entries.

This makes the payment traceable for reconciliation and reporting.

---

## Tech Stack

### Backend

- Go
- Gin
- GORM
- PostgreSQL
- JWT
- RBAC
- Wire Dependency Injection

### Documentation & API Testing

- Swagger / Swaggo
- Bruno REST Client
- PlantUML
- Markdown documentation

### Infrastructure

- Docker
- Docker Compose
- Environment-based configuration

---

## Architecture

The project follows a layered backend structure:

```text
Handler
  ↓
Service
  ↓
Repository
  ↓
Database
```

General responsibilities:

### Handler

- HTTP request binding
- Path/query/body parsing
- HTTP response formatting

### Service

- Business rules
- Validation
- Tenant isolation
- Transaction orchestration
- Cross-module coordination

### Repository

- Database queries
- Persistence operations
- GORM implementation details

### Domain / Model

- DTO contracts
- Service interfaces
- Repository interfaces
- Business models

---

## Database Domains

AkadiaPay separates the system into major domains:

### Master Domain

Responsible for identity and school data.

Includes:

- Tenant
- User
- Role
- UserTenantRole
- Student
- ParentStudent

### Payment Domain

Responsible for billing and payment operations.

Includes:

- PaymentPolicy
- PaymentProduct
- StudentObligation
- PaymentOrder
- PaymentAllocation
- LedgerEntry

---

## Demo Scenario

The seeded demo data supports a complete MVP story:

```text
School Admin logs in
  ↓
Creates payment policy
  ↓
Creates payment product
  ↓
Assigns product to student as obligation
  ↓
Parent logs in
  ↓
Parent views outstanding bills
  ↓
Parent creates payment order
  ↓
Payment is allocated to obligation
  ↓
Outstanding amount is reduced
  ↓
Ledger entry is posted
```

This demonstrates that AkadiaPay is not only a CRUD application.  
It models an actual school payment lifecycle.

---

## Documentation Included

This repository includes supporting documentation such as:

- MVP roadmap
- Phase-based implementation specifications
- Architecture review notes
- Bruno API collection
- PlantUML business-process diagrams
- Swagger/OpenAPI documentation

These documents were used to guide implementation and validate the overall system design.

---

## Repository Freeze Notice

This repository is intentionally frozen after the bootcamp presentation.

Reason:

Akadia is not intended to remain only as a bootcamp project.  
The next phase of Akadia will continue privately as a product initiative with commercial potential.

Future private development may include:

- QRIS integration
- Virtual Account integration
- Payment Gateway callback
- Settlement processing
- Notification engine
- WhatsApp and email reminder
- Advanced reporting
- Accounting module
- General ledger
- Student wallet
- Payment scheduling
- Financial analytics dashboard
- AI-based collection reminder

---

## Personal Note

AkadiaPay represents more than a final assignment.

It is the result of a focused learning journey in Golang backend development, product thinking, domain modeling, and AI-assisted software engineering.

The public MVP shows the foundation.

The real product journey continues privately.

---

## Author

Developed by **Norbertus Dewa Rucci**  
Built as part of Golang Backend Bootcamp final project, 2026.

---

## Disclaimer

This repository is provided as a public MVP snapshot for portfolio, learning, and demonstration purposes.

Production development, commercial features, and future business logic are not represented fully in this repository.