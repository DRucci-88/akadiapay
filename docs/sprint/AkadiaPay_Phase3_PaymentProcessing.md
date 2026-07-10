# AkadiaPay - Phase 3 (Payment Processing)

## Goal

Convert outstanding student obligations into completed payments while
preserving an auditable payment history.

## Business Flow

Parent ↓ View Outstanding Obligations ↓ Select Obligations ↓ Create
Payment Order ↓ Validate ↓ Allocate Payment ↓ Update Remaining Balance ↓
Close Obligation (if fully paid) ↓ Ready for Ledger

## Module 1 - Payment Order

Endpoints - POST /payment-orders - GET /payment-orders - GET
/payment-orders/:id - POST /payment-orders/:id/cancel

Request - StudentID - PaymentMethod - PaymentDate - Notes

Validation - Parent owns student - Student belongs to tenant -
Outstanding bill exists - Amount \> 0

## Module 2 - Payment Allocation

Endpoints - POST /payment-orders/:id/allocate - GET
/payment-orders/:id/allocations

Validation - Allocation \> 0 - Allocation \<= remaining balance - Total
allocation \<= payment amount - No duplicate obligations

## Repository

PaymentOrderRepository - Create - FindByID - FindPaginate - UpdateStatus

PaymentAllocationRepository - CreateBatch - FindByPaymentOrderID

StudentObligationRepository - UpdatePaidAmount - CloseIfCompleted

## Transaction

-   Create allocations
-   Update obligations
-   Update payment order
-   Prepare ledger

All in one transaction.

## Done Criteria

-   Payment Order
-   Payment Allocation
-   Partial payment
-   Full payment
-   Remaining balance
-   Ready for Ledger
