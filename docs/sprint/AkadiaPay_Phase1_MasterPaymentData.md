# AkadiaPay MVP - Phase 1 (Master Payment Data)

## Goal

Configure how and what a school bills before generating student
obligations.

## Business Flow

School Admin ↓ Create Payment Policies (Billing Strategies) ↓ Create
Payment Products ↓ Ready for Student Obligations

## Payment Policy (How to Bill)

-   Full Payment
-   Partial Payment
-   Installment
-   Allow Partial
-   Minimum Amount
-   Minimum Percentage
-   Allow Over Payment
-   Auto Close Obligation

One Payment Policy can be reused by many Payment Products.

## Payment Product (What to Bill)

-   SPP
-   Registration Fee
-   Building Fee
-   Uniform
-   Book Package

Each Payment Product references one Payment Policy.

## Deliverables

-   Payment Policy CRUD
-   Payment Product CRUD
-   Pagination
-   Filtering
-   Tenant Isolation
-   Seeder

## Next Phase

Student Obligation
