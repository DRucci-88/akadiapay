# AkadiaPay MVP Roadmap (After Authentication)

## Phase 0 - Authentication

-   Authentication (JWT)
-   Multi-tenant Login
-   Profile Endpoint
-   RBAC Foundation
-   Seeders: Roles, Tenants, Users, UserTenantRoles, Students,
    ParentStudents

## Phase 1 - Master Payment Data

### Payment Product

-   CRUD
-   Pagination
-   Filtering
-   Tenant Isolation

### Payment Policy

-   CRUD
-   Allow Partial Payment
-   Minimum Amount
-   Late Fee Strategy
-   Grace Period

## Phase 2 - Student Billing

### Student Obligation

-   Assign obligation
-   Outstanding bills
-   Remaining balance
-   Paid amount
-   Outstanding amount

## Phase 3 - Payment

### Payment Order

-   Create payment
-   Payment history
-   Pending / Paid / Cancelled

### Payment Allocation

-   Allocate one payment to one or many obligations

## Phase 4 - Ledger

### Ledger Entry

-   Auto create ledger after payment

## Phase 5 - Dashboard

### School Admin

-   Students
-   Outstanding
-   Revenue

### Parent

-   Children
-   Outstanding
-   Payment History

### Student

-   Bills
-   Payment History

## Demo Flow

1.  Login
2.  Payment Products
3.  Student Obligations
4.  Create Payment
5.  Allocate Payment
6.  Outstanding becomes zero
7.  Ledger created

## Backlog

-   Mapper Layer
-   JWT Manager
-   Remove Handler Interfaces
-   Redis
-   Notification
-   Payment Gateway
-   QRIS
-   Virtual Account
