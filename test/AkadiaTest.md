You are working on the Go Gin project AkadiaPay.

Phase 1 status:
- Completed service-layer tests for Payment Policy and Payment Product.
- All current Phase 1 executable tests are stored under `test/phase1` as requested.
- Phase 1 coverage runner is available at `test/phase1/run-coverage.ps1`.
- Working command:
  - `powershell -ExecutionPolicy Bypass -File .\test\phase1\run-coverage.ps1`
  - Default coverage output: `test/phase1/coverage.out`
- Scope tracker remains in `test/` and the implemented files are:
  - `test/phase1/helpers_test.go`
  - `test/phase1/payment_policy_service_test.go`
  - `test/phase1/payment_product_service_test.go`
  - `test/phase1/update_map_test.go`
  - `test/phase1/run-coverage.ps1`

Phase 2 status:
- Completed student obligation service tests under `test/phase2`.
- Phase 2 coverage runner is available at `test/phase2/run-coverage.ps1`.
- Implemented files:
  - `test/phase2/helpers_test.go`
  - `test/phase2/student_obligation_service_test.go`
  - `test/phase2/run-coverage.ps1`

Phase 3 status:
- Completed payment order service tests under `test/phase3`.
- Phase 3 coverage runner is available at `test/phase3/run-coverage.ps1`.
- Implemented files:
  - `test/phase3/helpers_test.go`
  - `test/phase3/payment_order_service_test.go`
  - `test/phase3/run-coverage.ps1`

Phase 4 status:
- Completed payment allocation service tests under `test/phase4`.
- Phase 4 coverage runner is available at `test/phase4/run-coverage.ps1`.
- Implemented files:
  - `test/phase4/helpers_test.go`
  - `test/phase4/payment_allocation_service_test.go`
  - `test/phase4/run-coverage.ps1`

Phase 5 status:
- Completed ledger entry service tests under `test/phase5`.
- Phase 5 coverage runner is available at `test/phase5/run-coverage.ps1`.
- Implemented files:
  - `test/phase5/helpers_test.go`
  - `test/phase5/ledger_entry_service_test.go`
  - `test/phase5/run-coverage.ps1`

Phase 6 status:
- Completed shared pagination helper tests under `test/phase6`.
- Phase 6 coverage runner is available at `test/phase6/run-coverage.ps1`.
- Implemented files:
  - `test/phase6/shared_test.go`
  - `test/phase6/run-coverage.ps1`

Phase 7 status:
- Added a minimal handler slice for payment policy binding and ID parsing behavior under `test/phase7`.
- Phase 7 coverage runner is available at `test/phase7/run-coverage.ps1`.
- Implemented files:
  - `test/phase7/payment_policy_handler_test.go`
  - `test/phase7/run-coverage.ps1`

Goal:
Generate meaningful unit tests for the existing AkadiaPay backend, focused on MVP business rules and service-layer correctness.

Important constraints:
- Do not rewrite business logic unless a test reveals a compile error or obvious bug.
- Do not refactor architecture broadly.
- Preserve current Clean Architecture layering.
- Prefer service-layer unit tests first.
- Repository/database integration tests are optional and should only be added if the existing project already has test DB patterns.
- Use table-driven tests where possible.
- Use testify if already available or add it if needed:
  go get github.com/stretchr/testify
- Keep tests deterministic.
- Do not depend on random data.
- Do not depend on external services.
- Do not require Docker unless integration tests already require it.
- Use mocks/fakes for repositories.
- Run gofmt on all generated test files.

Target command after implementation:
go test ./... -cover

Minimum target:
- Project compiles.
- All tests pass.
- Coverage improves meaningfully.
- Prioritize important business rules over high line coverage.

Project business flow to protect with tests:
Login → Payment Policy → Payment Product → Student Obligation → Payment Order → Allocation → Ledger

Main modules to test:
1. Auth
2. Payment Policy
3. Payment Product
4. Student Obligation
5. Payment Order
6. Payment Allocation
7. Ledger Entry
8. Shared pagination/update helpers

Implementation approach:

PHASE 1 — Inspect existing code
1. Inspect these folders:
   - domain/
   - internal/auth/
   - internal/payment/service/
   - internal/payment/repository/
   - internal/shared/
   - model/
2. Identify service constructors and interfaces.
3. Identify existing errors in shared/domain.
4. Identify DTO names and response names.
5. Identify whether repository interfaces are already in domain.
6. Generate tests using the actual existing names, not invented names.

PHASE 2 — Add testing helpers
Create test helper files only if useful.

Recommended:
internal/payment/service/test_helpers_test.go

Include:
- uuid constants for tenant, user, student, product, policy, obligation, order
- helper AuthContext builders:
  - adminAuthContext()
  - parentAuthContext()
  - studentAuthContext()
- helper model builders:
  - makePaymentPolicy()
  - makePaymentProduct()
  - makeStudentObligation()
  - makePaymentOrder()
  - makePaymentAllocation()
- float comparison helper:
  - assertAmountEqual(t, expected, actual)

If repository interfaces are difficult to mock manually, create simple fake structs inside each *_test.go file. Keep them small and focused on the method being tested.

PHASE 3 — Payment Policy service tests

Create:
internal/payment/service/payment_policy_service_test.go

Test cases:
1. Create full payment policy succeeds.
2. Create partial payment policy succeeds when minimum amount or minimum percentage is valid.
3. Create partial payment policy fails when both minimum amount and minimum percentage are zero.
4. Create policy fails when minimum amount is negative.
5. Create policy fails when minimum percentage is negative.
6. Create policy fails when minimum percentage is above 100.
7. Create policy forces minimum amount and percentage to zero when AllowPartial is false.
8. Update preserves repository error and does not call FindByID when update fails.
9. Update with pointer DTO supports false values.
10. Update validates final merged state, not only incoming DTO.

Important assertion:
- When req.AllowPartial = pointer(false), the update map/model must not skip false value.
- If current implementation uses Updates(map[string]any), assert false is included.

PHASE 4 — Payment Product service tests

Create:
internal/payment/service/payment_product_service_test.go

Test cases:
1. Create product succeeds with valid PaymentPolicyID.
2. Create product fails if price <= 0.
3. Create product fails if payment policy does not belong to tenant.
4. Create product stores RevenueAccountCode and RevenueAccountName when provided.
5. Update product supports PATCH pointer fields.
6. Update product validates final merged state.
7. FindByID uses tenant isolation.

PHASE 5 — Student Obligation service tests

Create:
internal/payment/service/student_obligation_service_test.go

Test cases:
1. Create obligation succeeds for valid student and payment product.
2. Create obligation fails if student does not belong to tenant.
3. Create obligation fails if payment product does not belong to tenant.
4. Create obligation fails if amount <= 0.
5. Create obligation initializes paid amount to 0.
6. Create obligation initializes outstanding amount equal to original amount.
7. Create obligation initializes status as pending/unpaid according to existing enum.
8. Bulk create succeeds for multiple students/products.
9. Bulk create rejects invalid student/product.
10. Update obligation supports status changes only if business rules allow it.
11. Delete/cancel obligation rejects already paid/closed obligation if existing business rule supports it.

PHASE 6 — Payment Order service tests

Create:
internal/payment/service/payment_order_service_test.go

Test cases:
1. Create payment order succeeds for valid student with outstanding obligations.
2. Create payment order fails when student does not belong to tenant.
3. Create payment order fails when amount <= 0.
4. Create payment order fails when amount exceeds total outstanding, because MVP does not support deposits.
5. Create payment order fails when no outstanding bill exists.
6. Create payment order succeeds for allowed payment methods.
7. Create payment order fails for invalid payment method.
8. Parent can create order only for linked student.
9. Parent cannot create order for unrelated student.
10. Cancel pending payment order succeeds.
11. Cancel completed payment order fails.
12. Cancel ledger-posted payment order fails if LedgerPostedAt is set.

PHASE 7 — Payment Allocation service tests

Create:
internal/payment/service/payment_allocation_service_test.go

This is the most important business-rule test file.

Test cases:
1. Allocate full payment succeeds.
2. Allocate partial payment succeeds when policy AllowPartial = true.
3. Allocate partial payment fails when policy AllowPartial = false.
4. Allocation below MinimumAmount fails.
5. Allocation below MinimumPercentage fails.
6. Allocation amount <= 0 fails.
7. Allocation exceeds obligation outstanding fails.
8. Total allocation exceeds payment order amount fails.
9. Duplicate obligation in same request fails.
10. Allocation to obligation from another student fails.
11. Allocation to obligation from another tenant fails.
12. Allocation to closed obligation fails.
13. Allocation to cancelled obligation fails.
14. Full allocation closes obligation when AutoCloseObligation = true.
15. Partial allocation keeps obligation partial/open.
16. Allocation updates PaidAmount and OutstandingAmount correctly.
17. Allocation completes payment order when allocated amount equals order total.
18. Allocation runs validation inside transaction / uses locking repository method if implemented.

Concurrency-related test:
- If service has LockByID / LockByIDs methods, assert fake repository records that lock methods were called during allocation.
- The purpose is not to test PostgreSQL locking, but to enforce service flow.

PHASE 8 — Ledger service tests

Create:
internal/payment/service/ledger_entry_service_test.go

Test cases:
1. PostPayment succeeds for completed payment order.
2. PostPayment fails for pending payment order.
3. PostPayment fails for cancelled payment order.
4. PostPayment is idempotent if ledger already exists.
5. PostPayment fails if payment order has no allocation.
6. PostPayment fails if allocation total does not equal payment order total.
7. PostPayment creates one debit entry.
8. PostPayment creates one or more credit entries.
9. Total debit equals total credit.
10. Cash payment maps to cash account.
11. Bank transfer maps to bank account.
12. QRIS maps to QRIS clearing account.
13. Virtual account maps to VA clearing account.
14. Product revenue account is used when RevenueAccountCode and RevenueAccountName exist.
15. Unknown product account falls back to configured OTHER revenue account if current implementation supports fallback.
16. LedgerPostedAt is set after successful posting.
17. Cross-tenant posting is rejected.

PHASE 9 — Shared helper tests

Create:
internal/shared/pageable_test.go
internal/shared/page_test.go
internal/shared/update_map_test.go if helper exists

Test cases:
1. Pageable Normalize defaults page <= 0 to 1.
2. Pageable Normalize defaults size <= 0 to default size.
3. Pageable Normalize caps size > max size.
4. Offset calculation is correct.
5. NewPage calculates total pages.
6. NewPage sets previous page correctly.
7. NewPage sets next page correctly.
8. NewPageEmpty returns empty data and valid pagination.
9. UpdateMap includes pointer false bool.
10. UpdateMap includes pointer zero float.
11. UpdateMap skips nil pointers.

PHASE 10 — Handler tests only if easy

If service mocks are straightforward, add handler tests using httptest.

Create:
internal/payment/handler/payment_policy_handler_test.go

Minimum handler tests:
1. POST with invalid JSON returns 422.
2. GET with invalid UUID returns 422 or 400 according to existing behavior.
3. Protected route behavior is not needed here if middleware is tested separately.

Do not spend too many changes on handler tests unless service tests are complete.

Mocking rules:
- Prefer manual fake repositories over generated mocks.
- Keep fake state simple.
- For each test, only implement methods used by that service method.
- If an interface has many methods, embed a fake base that panics on unexpected calls:

type fakePaymentPolicyRepository struct {
    updateFn func(...) (...)
}

func (f *fakePaymentPolicyRepository) Update(...) (...) {
    if f.updateFn == nil {
        panic("unexpected call: Update")
    }
    return f.updateFn(...)
}

Assertions:
- Use require.NoError for setup.
- Use require.Error for failure cases.
- Use assert.Equal for values.
- Use assert.ErrorIs when comparing sentinel errors.
- Use assert.InEpsilon for float amounts if project still uses float64.

Expected errors:
Use existing project errors if available:
- shared.ErrPaymentPolicyNotFound
- shared.ErrPaymentProductNotFound
- shared.ErrStudentObligationNotFound
- shared.ErrPaymentOrderNotFound
- shared.ErrInvalidAllocation
- shared.ErrOverAllocation
- shared.ErrPaymentPolicyPartialNotAllowed
- shared.ErrPaymentBelowMinimum
- etc.

If error names differ, inspect shared/domain errors and use actual names.
Do not invent new error names unless needed to compile and consistent with existing style.

Coverage priority order:
1. Payment Allocation service
2. Ledger service
3. Payment Order service
4. Student Obligation service
5. Payment Policy service
6. Payment Product service
7. Shared helpers
8. Handlers

After implementation run:
gofmt -w .
go test ./... -cover

If tests fail:
- Fix test assumptions first.
- Only fix production code when the test reveals an actual bug.
- Do not remove important business-rule tests just to pass quickly.

Final response should include:
- Test files created
- Business rules covered
- Coverage percentage if available
- Any production bugs fixed
- Remaining untested areas
