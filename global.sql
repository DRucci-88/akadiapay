-- 1. TABEL SEKOLAH / INSTANSI (Multi-Tenant Anchor)
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_name VARCHAR(255) NOT NULL,
    address TEXT,
    phone_number VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. TABEL PENGGUNA GLOBAL (Authentication & Profile Core)
-- Digunakan secara global oleh Orang Tua, Bendahara, Guru, maupun Admin Sistem
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE,
    phone_number VARCHAR(20) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. TABEL SISWA (Jangkar Utama untuk AkadiaPay, Learn, dan Report)
CREATE TABLE students (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    nis VARCHAR(50) NOT NULL, -- Nomor Induk Siswa
    full_name VARCHAR(255) NOT NULL,
    current_class VARCHAR(50), -- Struktur detail kelas bisa ditaruh di AkadiaReport/Learn nantinya
    status VARCHAR(20) DEFAULT 'ACTIVE', -- ACTIVE, GRADUATED, MUTATED, INACTIVE
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_nis_per_school UNIQUE (tenant_id, nis) -- NIS hanya unik di lingkup sekolah yang sama
);

-- 4. TABEL RELASI ORANG TUA - MURID (Many-to-Many Junction)
-- Mengakomodasi skenario 1 orang tua memiliki >1 anak di sekolah yang sama atau berbeda lintas tenant
CREATE TABLE parent_students (
    parent_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    relationship_type VARCHAR(50) DEFAULT 'FATHER', -- FATHER, MOTHER, GUARDIAN
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (parent_id, student_id)
);

-- 1. TABEL DAFTAR ROLE GLOBAL (Static/Master Data)
CREATE TABLE roles (
    id VARCHAR(50) PRIMARY KEY, -- 'SUPER_ADMIN', 'SCHOOL_ADMIN', 'TREASURER', 'TEACHER', 'PARENT', 'STUDENT'
    role_name VARCHAR(100) NOT NULL,
    description TEXT
);

-- 2. TABEL MAPPING USER TENANT ROLE (Jantung Kontrol Akses Akadia)
-- Mengatur peran spesifik seorang user di sekolah tertentu
CREATE TABLE user_tenant_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    role_id VARCHAR(50) NOT NULL REFERENCES roles(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Mencegah duplikasi role yang sama untuk user di satu sekolah
    CONSTRAINT unique_user_tenant_role UNIQUE (user_id, tenant_id, role_id)
)

-- Mempercepat query penarikan data siswa yang dimiliki oleh satu akun orang tua di Flutter
CREATE INDEX idx_parent_students_parent ON parent_students(parent_id);

-- Mempercepat query rule validation & ringkasan tagihan siswa aktif berdasarkan tenant di Go Backend
CREATE INDEX idx_obligations_student_status ON student_obligations(student_id, status);
CREATE INDEX idx_obligations_tenant ON student_obligations(tenant_id);

-- Mempercepat filter data ledger untuk kebutuhan riwayat keuangan/audit berdasarkan ID kewajiban
CREATE INDEX idx_ledgers_obligation ON financial_ledgers(obligation_id);