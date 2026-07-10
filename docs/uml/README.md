# Akadia Business Process PlantUML Pack

This directory contains PlantUML activity diagrams for the Akadia / School Financial Platform business processes.

Source basis:

- `docs/School Financial Platform Overview.docx`
- `docs/sprint/AkadiaPay_Phase0_Authentication.md`
- `docs/sprint/AkadiaPay_Phase1_MasterPaymentData.md`
- `docs/sprint/AkadiaPay_Phase2_StudentBilling.md`
- `docs/sprint/AkadiaPay_Phase3_PaymentProcessing.md`
- `docs/sprint/AkadiaPay_Phase4_FinancialLedger.md`
- Current source code routes and services in the latest project ZIP

## File Count

- 1 overview diagram: `BP-00`
- 28 detailed business process diagrams: `BP-01` to `BP-28`

## Current Code Coverage

Some diagrams describe implemented MVP behavior, while others document future scope from the source-of-truth document.

### Implemented or partially implemented in current MVP code

| Diagram | Current coverage |
|---|---|
| `BP-02_User_Login_and_Tenant_Workspace_Selection.puml` | Auth login/profile |
| `BP-08_Payment_Product_Configuration.puml` | PaymentProduct CRUD |
| `BP-09_Payment_Rule_and_Policy_Configuration.puml` | PaymentPolicy CRUD |
| `BP-15_Single_Student_Bill_Generation.puml` | StudentObligation create |
| `BP-16_Bulk_Bill_Generation.puml` | StudentObligation bulk create |
| `BP-18_Bill_Adjustment.puml` | StudentObligation update |
| `BP-19_Bill_Cancellation.puml` | StudentObligation delete/cancel behavior |
| `BP-21_Parent_Views_Multi_Student_Outstanding_Bills.puml` | Outstanding bills |
| `BP-22_Parent_Selects_Bills_and_Creates_Checkout.puml` | PaymentOrder create |
| `BP-26_Cash_Payment_through_Finance_Officer.puml` | Cash payment via PaymentOrder CASH |
| `BP-27_Payment_Allocation_Settlement_and_Financial_Ledger_Posting.puml` | Allocation, settlement and ledger posting |
| `BP-28_Notification_Reconciliation_and_Financial_Reporting.puml` | Ledger reporting basis |


### Future / roadmap processes from the source-of-truth document

These are included because the DOCX defines them as part of the target platform scope:

- Academic Year and Semester Setup
- Grade and Class Setup
- Dedicated Invoice module
- External Payment Gateway initiation
- Payment Gateway callback
- Scholarship
- Discount
- Penalty and Grace Period
- Notification Engine
- Full dashboard and reconciliation reports

## Recommended VS Code usage

1. Install PlantUML extension.
2. Open one `.puml` file.
3. Preview diagram.
4. Commit each business process separately if you want clean review history.

## Naming Convention

`BP-XX_Process_Name.puml`

`BP` means Business Process.
