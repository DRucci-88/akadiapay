# AkadiaPay MVP - Phase 0 (Authentication)

## Goal

Build a secure authentication and multi-tenant foundation before
implementing payment modules.

## Business Flow

User ↓ Login (Email + Password) ↓ Validate Credentials ↓ Find
UserTenantRoles ↓ Generate JWT ↓ Access Protected APIs

## Portals

-   Super Admin
-   School Admin
-   Parent
-   Student

## Business Rules

-   Email is globally unique.
-   JWT contains UserID, TenantID, RoleCode, optional StudentID.
-   TenantID always comes from authenticated context.
-   Every query is tenant isolated.

## Deliverables

-   Login
-   Profile
-   JWT Middleware
-   RBAC
-   Seeders

## Next Phase

Phase 1 - Master Payment Data
