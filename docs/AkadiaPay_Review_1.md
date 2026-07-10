# AkadiaPay Architecture Review & Technical Action Items
_Date: Current Review_

> Purpose:
>
> This document summarizes all architecture and implementation findings after reviewing
> Phases 0–4. It is intended to be used as implementation guidance for AI coding agents
> (Codex) and future contributors.

---

# Priority Legend

| Priority | Meaning |
|-----------|---------|
| P0 | Must fix before demo / correctness issue |
| P1 | Important business rule |
| P2 | Architecture improvement |
| P3 | Post-MVP enhancement |

---

# P0-01 Tenant Isolation on Detail Endpoints

## Problem

Repository methods such as `FirstByID()` query only by primary key.

```go
WHERE id = ?
```

This allows UUID enumeration across tenants.

## Required Design

Every payment entity lookup must include:

```sql
WHERE id = ?
AND tenant_id = ?
```

TenantID must always come from `AuthContext`.

Never trust client supplied tenant information.

## Affected Modules

- PaymentPolicy
- PaymentProduct
- StudentObligation
- PaymentOrder
- LedgerEntry

---

# P0-02 Preserve Repository Errors

Do not overwrite repository errors.

Bad:

```go
_, err := repo.Update(...)
entity, err := service.FindByID(...)
```

Correct:

```go
_, err := repo.Update(...)
if err != nil {
    return nil, err
}

return service.FindByID(...)
```

---

# P0-03 Concurrency Safe Allocation

## Problem

Current allocation validation occurs before transaction.

Concurrent requests can over allocate.

Example

Outstanding = 500000

Request A reads 500000

Request B reads 500000

Both allocate 400000

Result = 800000 allocated.

## Required Flow

Inside ONE transaction

1. Begin transaction
2. Lock PaymentOrder
3. Lock StudentObligation rows
4. Reload allocations
5. Validate balances again
6. Insert allocations
7. Update obligations
8. Commit

Use PostgreSQL

```go
clause.Locking{Strength:"UPDATE"}
```

---

# P0-04 RBAC Middleware

JWT authentication alone is insufficient.

Restrict routes.

Example

School Admin

- Payment Policy CRUD
- Payment Product CRUD
- Student Obligation CRUD
- Ledger

Parent

- Own students
- Own payments

Student

- Own obligations

Implement role middleware after JWT middleware.

---

# P1-01 Reject Unsupported Overpayment

Current MVP does not support student deposits.

Validation

```
PaymentAmount <= OutstandingAmount
```

Reject otherwise.

Future enhancement

```
Student Deposit Liability
```

---

# P1-02 Enforce Payment Policy

Payment Policy must affect allocation.

Rules

AllowPartial

```
false

allocation must equal outstanding
```

MinimumAmount

```
allocation >= minimum amount
```

MinimumPercentage

```
allocation >= outstanding * percentage
```

AllowOverPayment

```
false

reject excess allocation
```

AutoCloseObligation

```
true

status becomes CLOSED automatically
```

Never ignore Payment Policy during allocation.

---

# P1-03 PATCH Merge Validation

Update DTO uses pointer fields.

Validation must evaluate FINAL state.

Flow

1. Load existing entity.
2. Merge DTO into copy.
3. Validate merged copy.
4. Persist only supplied fields.

Never validate only incoming DTO.

---

# P2-01 Monetary Precision

Current

```
float64
```

Acceptable for MVP.

Recommended after MVP

- int64 (rupiah)
- decimal.Decimal

Never use equality comparison.

If float remains

```
abs(a-b) < epsilon
```

---

# P2-02 Ledger Mapping

Current implementation maps revenue accounts using product names.

Avoid

```
strings.Contains(name,"SPP")
```

Recommended

PaymentProduct contains

```
RevenueAccountCode
RevenueAccountName
```

or mapping table.

---

# P2-03 Ledger Meaning

Ledger entries represent accounting.

Debit

Increase Asset (Cash/Bank)

Credit

Increase Revenue

Do NOT describe them as increasing or decreasing student obligation.

Obligation tracking belongs to StudentObligation.

---

# P2-04 Ledger Immutability

Ledger history should never be modified.

Recommendations

- No update endpoint
- No delete endpoint
- Repository rejects modifications
- Consider removing soft delete later

---

# P2-05 Safe Sorting

Whitelist sortable columns.

Never concatenate client values directly into SQL.

Example

```go
allowed := map[string]string{
    "created_at":"created_at",
    "name":"name",
}
```

---

# P2-06 REST Naming

Prefer plural resources.

Good

```
/payment-policies
/payment-products
/payment-orders
/student-obligations
/ledger-entries
```

---

# General Service Rules

Service layer responsibilities

- Business validation
- Tenant validation
- Cross-schema composition
- Transaction orchestration
- Mapper invocation

Repository responsibilities

- Database only
- Query
- Insert
- Update
- Delete

Repositories must never contain business rules.

---

# Transaction Rules

Transactions should wrap business operations, not repositories.

Example

```
PaymentOrder
    ↓
PaymentAllocation
    ↓
StudentObligation
    ↓
LedgerPosting
```

Everything above executes inside one transaction.

---

# Security Rules

Always derive

- TenantID
- UserID
- StudentID
- RoleCode

from JWT.

Never accept them from request body.

---

# Coding Standards

- Response DTO != Entity
- PATCH uses pointer DTO + Updates(map[string]any)
- Repository returns entities
- Service returns DTOs
- Handler only orchestrates HTTP

---

# Acceptance Criteria

Before Phase 5:

- Tenant isolation verified
- RBAC enforced
- Allocation concurrency safe
- Payment Policy enforced
- Ledger idempotent
- Overpayment rejected
- PATCH merge validation implemented
- Integration tests for payment flow
