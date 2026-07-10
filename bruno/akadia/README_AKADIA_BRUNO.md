# Akadia Bruno API Collection

This collection is designed for the seeded Akadia MVP dataset.

## Required manual setup

1. Run the backend and database seeder.
2. Open the `dev_akadia` environment.
3. Run the login requests under `auth/`.
4. Copy the returned token for each role into the matching environment variable.
5. Run the reference requests under `00 Reference - Seeded IDs/`.
6. Copy IDs from responses into the environment variables.

## Important accounts

All passwords are:

```text
password
```

Recommended tokens:

| Variable | Login account | Role / tenant |
|---|---|---|
| token_platform | admin@akadia.id | SUPER_ADMIN / AKADIA |
| token_school | admin@sman1.id | SCHOOL_ADMIN / SMAN1 |
| token_treasurer | treasurer@sman1.id | TREASURER / SMAN1 |
| token_parent | budi.parent@gmail.com | PARENT / SMAN1 |
| token_student | kevin@student.id | STUDENT / SMAN1 |
| token_school_smahb | admin@harapan.id | SCHOOL_ADMIN / SMAHB |
| token_treasurer_smahb | finance@harapan.id | TREASURER / SMAHB |
| token_parent_smahb | asep.parent@gmail.com | PARENT / SMAHB |
| token_student_smahb | gilis@student.id | STUDENT / SMAHB |

## Seeded story

SMAN1:

- Kevin has closed SPP July, open SPP August, closed Book Package, partial Study Tour, and cancelled Uniform.
- Rucco has partial Building Fee and pending SPP August.
- Budi is parent for Kevin and Rucco.

SMAHB:

- Gilis has closed SPP July, partial Building Fee, pending Study Tour, and cancelled Book Package.
- Asep is parent for Gilis.

## Recommended workflow

1. `auth` - generate tokens.
2. `00 Reference - Seeded IDs` - collect IDs.
3. `01 Payment Policies` - test configuration.
4. `02 Payment Products` - test products and account mapping.
5. `03 Student Obligations` - test billing.
6. `04 Parent Outstanding` - test parent/student-facing outstanding bills.
7. `05 Payment Orders` - test checkout/payment order.
8. `06 Payment Allocations` - test allocation, partial/full payment, and ledger auto-posting.
9. `07 Ledger Entries` - test financial ledger queries.
10. `90 Edge Cases` - test expected failures.

## Notes

This collection uses environment variables for IDs because database UUIDs are generated during seeding. The seed data has stable business codes, but not stable UUIDs.
