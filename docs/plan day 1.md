# AkadiaPay - Sprint Day 1

> Goal: Build the foundation of AkadiaPay.
>
> ❌ No payment business logic today.
> ✅ Finish architecture, authentication, and Master Module.

---

# Today's Deliverables

- [ ] Project architecture finalized
- [ ] PostgreSQL schemas created
- [ ] GORM CLI generated
- [ ] Authentication working
- [ ] RBAC working
- [ ] Master Module CRUD
- [ ] Seeder
- [ ] Bruno Collection
- [ ] Commit Day 1

---

# 1. Freeze Architecture

Project Structure

```
cmd/
app/
internal/
    auth/
    master/
        dto/
        handler/
        repository/
        service/
    payment/
        dto/
        handler/
        repository/
        service/
    shared/
        config/
        helper/
        validator/
        middleware/
        response/
        pagination/
model/
    generated/
migration/
seeder/
docs/
bruno/
docker/
```

---

# 2. Database

Create only two PostgreSQL schemas.

```
master
payment
```

The SQL files from the business design become the reference.

Do NOT create additional schemas today.

---

# 3. Domain Models

## Schema master

### User

Authentication account.

Role:
- SUPER_ADMIN
- SCHOOL_ADMIN
- PARENT

---

### School

Represents one institution.

Examples:

- SMK Negeri 1
- SMP Harapan Bangsa

CRUD Required.

---

### Student

Represents one student.

Relationship

Student

↓

School

CRUD Required.

---

### PaymentProduct

Examples

- SPP
- Study Tour
- Registration
- Books
- Laboratory

CRUD Required.

This is only a TEMPLATE.

Not a bill.

---

## Schema payment

Only prepare models.

NO BUSINESS LOGIC.

Create

- StudentObligation
- Invoice
- Payment
- LedgerEntry
- PaymentPolicy

Repository/service can remain empty today.

---

# 4. GORM CLI

After models are finished

Generate

- generated query
- metadata
- field maps

Verify

- AutoMigrate
- Relationships
- Foreign Keys

---

# 5. Authentication

Implement

- Login
- JWT
- Middleware
- Authorization

RBAC

- SUPER_ADMIN
- SCHOOL_ADMIN
- PARENT

No Refresh Token yet.

---

# 6. CRUD

Implement CRUD only for

## School

- Create
- Update
- Delete
- FindByID
- FindAll
- Pagination
- Search

---

## Student

- Create
- Update
- Delete
- FindByID
- FindAll
- Pagination
- Search

---

## Payment Product

- Create
- Update
- Delete
- FindByID
- FindAll
- Pagination
- Search

---

# 7. Seeder

Generate demo data

School

↓

Students

↓

Payment Products

Do NOT seed

- Obligations
- Invoice
- Payment
- Ledger

Those belong to Day 2.

---

# 8. Bruno

Prepare APIs

POST /login

School

POST
GET
PUT
DELETE

Student

POST
GET
PUT
DELETE

Payment Product

POST
GET
PUT
DELETE

---

# 9. Docker

Only verify

```
docker compose up
```

Application starts successfully.

Database connected.

AutoMigrate works.

Seeder works.

---

# 10. Commit

Commit message

```
feat(day1): bootstrap master module and authentication
```

---

# Out of Scope Today

Do NOT touch

- Student Obligation
- Invoice
- Payment
- Ledger
- Payment Calculation
- Payment Gateway
- QRIS
- Notification
- Reports
- Dashboard

Those belong to Day 2 and Day 3.

---

# Success Criteria

At the end of Day 1 we should have

✅ Authentication

✅ RBAC

✅ School CRUD

✅ Student CRUD

✅ Payment Product CRUD

✅ Seeder

✅ Docker

✅ Bruno Collection

✅ Clean Architecture

If all of the above are completed, Day 2 will focus entirely on business logic.

---

# Important Architecture Decisions

Decision #001

Model IS Domain.

No separate `domain/` folder.

---

Decision #002

Monolith First.

One Go application.

Two PostgreSQL schemas.

Future modules may become independent services without rewriting the business logic.

---

Decision #003

Everything revolves around Student Obligation.

Payments are events.

Student Obligation is state.

Never model Payment as the center of the system.