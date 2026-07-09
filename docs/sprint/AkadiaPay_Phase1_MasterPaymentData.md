# AkadiaPay MVP - Phase 1 (Updated)

## Goal

Build the payment configuration layer before any billing can occur.

## Business Flow

School Admin ↓ Create Payment Policies (Billing Strategies) ↓ Create
Payment Products ↓ Assign Products to Students

------------------------------------------------------------------------

## 1. Payment Policy (Billing Strategy)

Payment Policy defines **how a payment product can be paid**.

Examples: - Full Payment - Partial Payment - Installment

### CRUD

-   Create
-   Update
-   Delete
-   Detail
-   List

### Required Fields

-   TenantID
-   Code
-   Name
-   Description
-   AllowPartial
-   MinimumAmount
-   MinimumPercentage
-   AllowOverPayment
-   AutoCloseObligation

### Validation

-   Code unique per tenant
-   Name required
-   MinimumAmount / MinimumPercentage only used when partial payment is
    enabled

### Seeder

-   FULL_PAYMENT
-   PARTIAL_PAYMENT
-   INSTALLMENT

------------------------------------------------------------------------

## 2. Payment Product

Payment Product defines **what the school bills**.

Examples: - SPP - Registration Fee - Building Fee - Uniform - Book
Package

Each Payment Product references one Payment Policy.

### CRUD

-   Create
-   Update
-   Delete
-   Detail
-   List

### Required Fields

-   TenantID
-   PaymentPolicyID
-   Code
-   Name
-   Description
-   Category
-   Amount
-   IsActive

### Validation

-   Code unique per tenant
-   Amount \> 0
-   PaymentPolicy must exist

### Seeder

-   SPP
-   Registration
-   Uniform
-   Building Fee
-   Book Package

------------------------------------------------------------------------

## Deliverables

-   Repository
-   Service
-   Handler
-   REST Client
-   Seeder
-   Pagination
-   Filtering
-   Tenant Isolation

------------------------------------------------------------------------

## Done Criteria

-   Payment Policy CRUD
-   Payment Product CRUD
-   Pagination
-   Filtering
-   Tenant Isolation
-   Seeder

------------------------------------------------------------------------

## Next Phase

Student Obligation

Assign Payment Products to Students to generate bills.
