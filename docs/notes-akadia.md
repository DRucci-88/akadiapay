# Hackaton Akadia

Akadia
- AkadiaPay
- AkadiaPayroll
- AkadiaInventory
- AkadiaLedger
- AkadiaCourse

Modular Monolith
Domain Driven Design Architecture

## Packages

```txt
go get github.com/gin-gonic/gin ;
go get gorm.io/gorm ;
go get gorm.io/driver/postgres ;
go get -u golang.org/x/crypto/bcrypt ;
go get -u github.com/golang-jwt/jwt/v5 ;
go get github.com/google/wire ;
go get -u github.com/swaggo/gin-swagger ;
go get -u github.com/swaggo/files ;
go get golang.org/x/sync/errgroup ;
go get github.com/brianvoe/gofakeit/v7 ;
go get github.com/joho/godotenv ;
go get gorm.io/cli/gorm/field ;
go get gorm.io/cli/gorm/genconfig ;


go install github.com/google/wire/cmd/wire@latest
go install github.com/swaggo/swag/cmd/swag@latest
```

## Bootcamp Requirement
- Swagger Wajib
- Authentication & Authorization (jwt + RBAC) 
- CRUD Lengkap (PostgreeSQL + Gorm) 
- Pagination & Flitering 
- Upload file optional 
- Docker Deployment 
- Clean Architecture 
- Minimal 60% Coverage

## Mental Roadmap

```txt
Pembayaran bukan objek utama. Objek utama adalah kewajiban yang harus diselesaikan oleh siswa.
```

Deliverables:

MVP Scope
Domain Model
Bounded Context
Folder Structure
API List
Database Review

## MVP 

```txt
A parent can pay a school bill online for their children.
```

### MVP Story
```txt
Admin Creates -> Payment Product [SPP July 500k]
Admin Assign John -> Creates Student Obligation
Parent Opens Outstanding Biils -> Choose SPP July
Parent Pays 250k -> Remaining 250k
Parent pays again 250l -> Completed
```

## Mental Mode

```txt
Student Obligation is the core feature of AkadiaPay

Payments are events. Obligations are state.
```

### Domain Model Scopes
master
- School
- Student
- PaymentProduct

payment
- StudentObligation
- PaymentRule / PaymentPolicy
- Invoice
- Payment
- Ledger

## Sprint 4 Days


- Architecture & schemas (master, payment)
- Authentication & RBAC
- School, Student, Payment Product CRUD
- Seeder


- Student Obligation
- Payment Policy (minimum/full/partial)
- Obligation assignment
- Outstanding bill listing


- Invoice creation
- Payment flow (simulated)
- Ledger posting
- Remaining balance calculation


- Docker
- Tests (focus on payment policy and obligation logic to help reach the coverage target)
- README
- Demo data
- Presentation
- Polish


## Project Structure
```txt
cmd/
    api/
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
        pagination/
        response/
        errors/
model/
    generated/
migration/
seeder/
docs/
bruno/
docker/
scripts/
test/
```

## Handler

```txt
POST /login
POST /schools
GET /schools
PUT /schools
DELETE /schools
POST /students
GET /students
POST /payment-products
GET /payment-products
```

## Flow Development
domain
repository
service
dto factory
handler