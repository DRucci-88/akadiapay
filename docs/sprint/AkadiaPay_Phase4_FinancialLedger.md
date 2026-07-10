# AkadiaPay — Phase 4: Financial Ledger

> Purpose  
> Record every completed payment as immutable accounting evidence so AkadiaPay can support reconciliation, reporting, audits, and future AkadiaLedger integration.

---

# 1. Scope

Phase 4 begins **after a Payment Order has been successfully completed**.

It must:

- Create balanced debit and credit ledger entries.
- Link ledger entries to the Payment Order.
- Prevent duplicate ledger posting.
- Preserve immutable financial history.
- Support ledger listing and payment-based lookup.
- Keep posting inside the same transaction as payment completion whenever possible.

Phase 4 does **not** implement:

- Full Chart of Accounts management.
- General journal UI.
- Manual journal adjustments.
- Reversal workflow.
- Closing periods.
- Tax accounting.
- External accounting integration.

Those belong to future versions.

---

# 2. Main Business Concepts

## 2.1 Payment Order

Represents the payment transaction.

Example:

```text
Payment Order: PAY-20260712-0001
Student: Kevin Wijaya
Amount: Rp700.000
Status: COMPLETED
```

## 2.2 Payment Allocation

Explains how the payment was distributed.

Example:

```text
SPP July       Rp500.000
Book Package   Rp200.000
```

## 2.3 Ledger Entry

Represents one debit or credit line.

Example:

```text
Debit  Cash / Bank                Rp700.000
Credit Tuition Revenue            Rp500.000
Credit Book Revenue               Rp200.000
```

The sum of debit must always equal the sum of credit.

---

# 3. Aggregate and Ownership

Primary aggregate for this phase:

```text
PaymentOrder
    └── LedgerEntries
```

Posting is triggered from the Payment Order completion flow.

A Ledger Entry must never be created independently by an API client in the MVP.

---

# 4. Business Flow

## 4.1 Successful Payment Posting

```text
Payment Order is COMPLETED
        │
        ▼
Check ledger has not been posted
        │
        ▼
Determine debit account
        │
        ▼
Determine one or more credit accounts
        │
        ▼
Build ledger entries
        │
        ▼
Validate debit = credit
        │
        ▼
Insert ledger entries
        │
        ▼
Mark Payment Order as ledger-posted
        │
        ▼
Commit transaction
```

## 4.2 Example

```text
Parent pays Rp700.000

Allocation:
- SPP July: Rp500.000
- Books: Rp200.000

Ledger:
- Debit Cash: Rp700.000
- Credit Tuition Revenue: Rp500.000
- Credit Book Revenue: Rp200.000
```

---

# 5. Technical Flow

## 5.1 Recommended Transaction Flow

Everything below should execute in one database transaction:

1. Lock Payment Order.
2. Verify Payment Order status is `COMPLETED`.
3. Verify ledger has not already been posted.
4. Load Payment Allocations.
5. Load related Payment Products.
6. Resolve debit account based on payment method.
7. Resolve credit account for each product.
8. Build Ledger Entry models.
9. Validate totals.
10. Insert Ledger Entries in batch.
11. Update Payment Order ledger status.
12. Commit.

Pseudo flow:

```go
repo.Transaction(ctx, func(txRepo domain.RepositoryManagerPayment) error {
	order := txRepo.PaymentOrder().LockByID(...)
	validateOrder(order)

	allocations := txRepo.PaymentAllocation().FindByOrderID(...)
	entries := ledgerBuilder.Build(order, allocations)

	validateBalanced(entries)

	txRepo.LedgerEntry().CreateInBatches(...)
	txRepo.PaymentOrder().MarkLedgerPosted(...)

	return nil
})
```

---

# 6. Data Model

Recommended MVP model:

```go
type LedgerEntry struct {
	BaseModel

	PaymentOrderID uuid.UUID `gorm:"type:uuid;not null;index:idx_ledger_entry_order"`
	PaymentOrder   *PaymentOrder `gorm:"foreignKey:PaymentOrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	EntryDate time.Time `gorm:"type:date;not null;index"`

	AccountCode string `gorm:"type:varchar(30);not null;index"`
	AccountName string `gorm:"type:varchar(150);not null"`

	Debit  float64 `gorm:"type:numeric(18,2);default:0;not null"`
	Credit float64 `gorm:"type:numeric(18,2);default:0;not null"`

	Description string `gorm:"type:text"`
}
```

Recommended addition to `PaymentOrder`:

```go
LedgerPostedAt *time.Time `gorm:"index"`
```

Optional:

```go
LedgerPostingStatus LedgerPostingStatus
```

Possible values:

```text
NOT_POSTED
POSTED
FAILED
```

For MVP, `LedgerPostedAt != nil` is enough.

---

# 7. Account Mapping

For MVP, use fixed account codes.

## 7.1 Debit Account by Payment Method

```text
CASH             → 1101 Cash
BANK_TRANSFER    → 1102 Bank
QRIS             → 1103 QRIS Clearing
VIRTUAL_ACCOUNT  → 1104 Virtual Account Clearing
```

## 7.2 Credit Account by Payment Product Category

```text
TUITION       → 4101 Tuition Revenue
REGISTRATION  → 4102 Registration Revenue
BUILDING      → 4103 Building Revenue
UNIFORM       → 4104 Uniform Revenue
BOOK          → 4105 Book Revenue
OTHER         → 4199 Other Education Revenue
```

The mapping can initially live in service-level constants.

Future version:

```text
ChartOfAccounts
LedgerAccountMapping
```

---

# 8. Endpoints

## 8.1 Get Ledger Entries by Payment Order

```http
GET /payment-orders/:id/ledger
```

Purpose:

- Show accounting entries created for one payment.
- Used by school admin and audit screens.

Validation:

- Payment Order must belong to authenticated tenant.
- Payment Order must exist.
- Caller must have permitted role.

## 8.2 List Ledger Entries

```http
GET /ledger-entries
```

Filters:

- `payment_order_id`
- `entry_date_from`
- `entry_date_to`
- `account_code`
- `keyword`
- `page`
- `size`

Tenant isolation is mandatory.

## 8.3 Manual Posting Endpoint

Recommended for fallback only:

```http
POST /payment-orders/:id/post-ledger
```

Use cases:

- Previous automatic posting failed.
- Admin retries posting.
- Migrated payment order has not been posted.

This endpoint must be idempotent.

---

# 9. Service Contracts

## 9.1 Ledger Service

```go
type LedgerService interface {
	PostPayment(
		ctx context.Context,
		authContext *security.AuthContext,
		paymentOrderID uuid.UUID,
	) error

	FindByPaymentOrderID(
		ctx context.Context,
		authContext *security.AuthContext,
		paymentOrderID uuid.UUID,
	) ([]LedgerEntryResponse, error)

	FindPaginate(
		ctx context.Context,
		authContext *security.AuthContext,
		pageable *shared.Pageable,
		filter *LedgerEntryFilter,
	) (*shared.Page[LedgerEntryResponse], error)
}
```

## 9.2 Ledger Repository

```go
type LedgerEntryRepository interface {
	FindByPaymentOrderID(
		ctx context.Context,
		paymentOrderID uuid.UUID,
	) ([]model.LedgerEntry, error)

	ExistsByPaymentOrderID(
		ctx context.Context,
		paymentOrderID uuid.UUID,
	) (bool, error)

	CreateInBatches(
		ctx context.Context,
		entries []model.LedgerEntry,
		batchSize int,
	) error

	FindPaginate(
		ctx context.Context,
		tenantID uuid.UUID,
		pageable *shared.Pageable,
		filter *LedgerEntryFilter,
	) (*shared.Page[model.LedgerEntry], error)
}
```

---

# 10. Validation Rules

Validation belongs in the Service Layer.

## 10.1 Payment Order Validation

- Payment Order must exist.
- Payment Order must belong to authenticated tenant.
- Payment Order status must be `COMPLETED`.
- Cancelled Payment Order cannot be posted.
- Pending Payment Order cannot be posted.

## 10.2 Posting Validation

- Ledger must not already exist for the Payment Order.
- Payment Order must have at least one allocation.
- Total allocated amount must equal Payment Order total amount.
- Total debit must equal total credit.
- Every entry amount must be positive on exactly one side.
- A single entry cannot have both debit and credit greater than zero.
- A single entry cannot have both debit and credit equal to zero.

## 10.3 Account Mapping Validation

- Payment method must map to a debit account.
- Payment product category must map to a credit account.
- Unknown categories may use `4199 Other Education Revenue`.
- Missing debit mapping must fail the posting.

---

# 11. Edge Cases

## 11.1 Duplicate Posting

Expected:

- Do not insert duplicate entries.
- Return existing result or `ErrLedgerAlreadyPosted`.
- Recommended MVP behavior: return existing entries with HTTP 200.

## 11.2 Payment Order Completed but No Allocation

Expected:

- Reject posting.
- Return `ErrPaymentAllocationNotFound`.

## 11.3 Allocation Total Does Not Match Payment Total

Expected:

- Roll back transaction.
- Return `ErrLedgerUnbalancedSource`.

## 11.4 Unknown Payment Method

Expected:

- Reject posting.
- Return `ErrLedgerDebitAccountNotConfigured`.

## 11.5 Unknown Product Category

Expected:

- Use fallback revenue account `4199`.
- Log a warning.
- Do not fail unless strict mode is added later.

## 11.6 Zero Amount Allocation

Expected:

- Reject posting.
- Never generate a zero-value ledger entry.

## 11.7 Payment Cancellation After Posting

MVP rule:

- A posted payment cannot be deleted.
- Cancellation must not remove ledger history.
- Future version should create reversal entries.

Recommended current behavior:

```text
if LedgerPostedAt != nil:
    reject cancellation with ErrPostedPaymentCannotBeCancelled
```

## 11.8 Database Failure During Batch Insert

Expected:

- Entire transaction rolls back.
- Payment Order remains unposted.
- Retry is safe.

## 11.9 Concurrent Posting Requests

Expected:

- Lock Payment Order row.
- Recheck `LedgerPostedAt` after lock.
- Only one request inserts entries.
- Second request returns already-posted result.

## 11.10 Partial Payment

A partial payment still creates balanced entries for the amount actually paid.

```text
Debit Cash: Rp200.000
Credit Tuition Revenue: Rp200.000
```

## 11.11 Overpayment

Recommended MVP simplification:

```text
Do not support unallocated overpayment in Phase 4.
Require total payment = total allocation.
```

Future account:

```text
2101 Customer Credit / Student Deposit
```

## 11.12 Soft-Deleted Payment Product

Historical posting must still work.

Recommended:

- Retain snapshot fields in payment/allocation records.
- Do not depend on active-only product queries when posting historical data.

## 11.13 Tenant Isolation Failure

Every lookup must ensure:

```text
PaymentOrder.TenantID == AuthContext.TenantID
```

Never post a payment for another tenant.

---

# 12. Idempotency

Posting must be idempotent.

For MVP, use:

```text
PaymentOrder.LedgerPostedAt
+
ExistsByPaymentOrderID()
+
row lock
```

---

# 13. Error Definitions

```go
var (
	ErrLedgerAlreadyPosted              = errors.New("ledger already posted")
	ErrLedgerEntryNotFound              = errors.New("ledger entry not found")
	ErrLedgerUnbalanced                 = errors.New("ledger entries are not balanced")
	ErrLedgerUnbalancedSource           = errors.New("payment allocation does not match payment total")
	ErrLedgerDebitAccountNotConfigured  = errors.New("debit account is not configured")
	ErrLedgerCreditAccountNotConfigured = errors.New("credit account is not configured")
	ErrPaymentAllocationNotFound        = errors.New("payment allocation not found")
	ErrPostedPaymentCannotBeCancelled   = errors.New("posted payment cannot be cancelled")
)
```

---

# 14. DTOs

## Ledger Entry Response

```go
type LedgerEntryResponse struct {
	ID             uuid.UUID `json:"id"`
	PaymentOrderID uuid.UUID `json:"payment_order_id"`
	EntryDate      time.Time `json:"entry_date"`
	AccountCode    string    `json:"account_code"`
	AccountName    string    `json:"account_name"`
	Debit          float64   `json:"debit"`
	Credit         float64   `json:"credit"`
	Description    string    `json:"description"`
}
```

## Ledger Filter

```go
type LedgerEntryFilter struct {
	PaymentOrderID *uuid.UUID `form:"payment_order_id"`
	AccountCode    *string    `form:"account_code"`
	EntryDateFrom  *time.Time `form:"entry_date_from" time_format:"2006-01-02"`
	EntryDateTo    *time.Time `form:"entry_date_to" time_format:"2006-01-02"`
	Keyword        *string    `form:"keyword"`
}
```

---

# 15. Posting Algorithm

```text
Input: PaymentOrderID

1. Begin transaction.
2. Lock PaymentOrder.
3. Validate tenant.
4. Validate status = COMPLETED.
5. Check LedgerPostedAt.
6. Load allocations.
7. Sum allocations.
8. Compare with PaymentOrder.TotalAmount.
9. Resolve debit account.
10. Group allocations by revenue account.
11. Create one debit entry.
12. Create one credit entry per revenue account.
13. Validate debit total = credit total.
14. Insert entries in batch.
15. Set LedgerPostedAt = now.
16. Commit.
```

Grouping allocations avoids duplicate credit rows.

---

# 16. Security and Authorization

Allowed roles:

```text
SUPER_ADMIN
SCHOOL_ADMIN
```

For MVP:

- Parents and students view receipts.
- They do not access raw ledger endpoints.

---

# 17. Logging

Log:

- Payment Order ID
- Tenant ID
- Posting result
- Entry count
- Debit total
- Credit total
- Error

Never log:

- JWT token
- Password
- Sensitive payment credentials

---

# 18. Testing Requirements

Minimum unit tests:

1. Full payment creates balanced entries.
2. Partial payment creates balanced entries.
3. Duplicate posting is idempotent.
4. Pending payment rejected.
5. Cancelled payment rejected.
6. Allocation mismatch rejected.
7. Unknown payment method rejected.
8. Unknown category uses fallback.
9. Cross-tenant posting rejected.
10. Concurrent posting does not duplicate entries.

Integration test:

```text
Create obligation
→ create payment order
→ allocate payment
→ complete payment
→ post ledger
→ verify debit = credit
```

---

# 19. Acceptance Criteria

Phase 4 is complete when:

- Completed payments automatically create ledger entries.
- Debit equals credit for every payment.
- Duplicate posting is prevented.
- Ledger history is immutable.
- Ledger entries can be queried by Payment Order.
- Tenant isolation is enforced.
- Retry is safe.
- Integration test passes.

---

# 20. Future Enhancements

- Chart of Accounts.
- Configurable account mapping.
- Reversal entries.
- Refund accounting.
- Student deposit liability.
- General journal.
- Accounting periods.
- External accounting export.
- AkadiaLedger integration.
