# AkadiaPay MVP - Phase 2 (Student Billing)

## Goal

Generate billable obligations for students. At the end of this phase,
every student can have outstanding balances before any payment is made.

------------------------------------------------------------------------

# Business Flow

``` text
School Administrator
        │
        ▼
Select Student(s)
        │
        ▼
Select Payment Product
        │
        ▼
System copies Billing Strategy from Payment Policy
        │
        ▼
Create Student Obligation
        │
        ▼
Student / Parent can view Outstanding Bills
```

------------------------------------------------------------------------

# Core Concept

Payment Policy = How to bill

Payment Product = What to bill

Student Obligation = Who must pay

The Student Obligation is the center of the payment domain.

------------------------------------------------------------------------

# Endpoints

## POST /student-obligations

Create one obligation.

Request

-   StudentID
-   PaymentProductID
-   DueDate
-   Amount (optional, default from product)
-   Notes

Validation

Handler - UUID format - JSON validation

Service - Student exists - Student belongs to current tenant - Payment
Product exists - Payment Product belongs to tenant - Payment Policy
exists - No duplicate active obligation for same Student + Product +
Period (business rule) - Amount \> 0 - DueDate required

Repository - Insert obligation

Response - Created obligation

------------------------------------------------------------------------

## POST /student-obligations/bulk

Assign one payment product to many students.

Request

-   PaymentProductID
-   StudentIDs\[\]
-   DueDate

Validation

-   StudentIDs not empty
-   Ignore duplicated StudentIDs
-   Transaction required
-   Rollback entire batch on failure

------------------------------------------------------------------------

## GET /student-obligations

Filters

-   StudentID
-   PaymentProductID
-   Status
-   Keyword
-   DueDateFrom
-   DueDateTo

Support

-   Pagination
-   Sorting
-   Tenant Isolation

------------------------------------------------------------------------

## GET /student-obligations/{id}

Return obligation detail.

------------------------------------------------------------------------

## PUT /student-obligations/{id}

Editable fields

-   DueDate
-   Notes

Business Rule

Amount and PaymentProduct cannot be changed after payments exist.

------------------------------------------------------------------------

## DELETE /student-obligations/{id}

Soft delete only.

Cannot delete if payment allocations already exist.

------------------------------------------------------------------------

## GET /students/{studentId}/outstanding

Return

-   All unpaid obligations
-   Paid amount
-   Remaining amount
-   Total outstanding

Used by

-   Parent Portal
-   Student Portal
-   Payment screen

------------------------------------------------------------------------

# Suggested Status

PENDING

PARTIAL

PAID

CANCELLED

------------------------------------------------------------------------

# Technical Flow

Handler ↓

DTO Validation

↓

Service

Business Validation

↓

Repository

Transaction (if needed)

↓

Database

------------------------------------------------------------------------

# Service Responsibilities

-   Validate ownership
-   Calculate remaining amount
-   Prevent duplicate obligation
-   Prevent invalid updates
-   Tenant isolation

------------------------------------------------------------------------

# Repository Responsibilities

-   CRUD
-   Pagination
-   Aggregate outstanding amount
-   Bulk insert
-   Transaction support

------------------------------------------------------------------------

# Seeder

Generate

-   3 payment products and its payment policy
-   3 students
-   6 obligations

Example

Kevin - SPP July - Book Fee

Rucco - SPP July - Registration

Gilis - SPP July - Uniform

------------------------------------------------------------------------

# Done Criteria

-   CRUD works
-   Bulk assignment works
-   Outstanding endpoint works
-   Pagination works
-   Tenant isolation works

------------------------------------------------------------------------

# Next Phase

Payment Processing

Create Payment Orders and allocate payments to Student Obligations.
