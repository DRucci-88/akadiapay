-- 5. TABEL KEBIJAKAN ATURAN PEMBAYARAN (Payment Policy Table)
-- Memisahkan aturan umum dengan tipe data ketat dan menyediakan kolom dinamis
CREATE TABLE payment_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    policy_name VARCHAR(100) NOT NULL, -- Contoh: "Kebijakan SPP Bulanan", "Kebijakan Tour Minimal DP"
    
    -- Konfigurasi Parameter Dasar & Cicilan
    allow_installment BOOLEAN DEFAULT FALSE,       -- Apakah boleh dicicil / partial payment?
    min_payment_amount NUMERIC(12, 2) DEFAULT 0,   -- Batas minimum rupiah sekali bayar (jika dicicil)
    min_payment_percentage NUMERIC(5, 2) DEFAULT 0, -- Batas minimum persentase dari total tagihan (misal: DP 30%)
    
    -- Konfigurasi Parameter Denda & Tenggang Waktu
    has_late_fee BOOLEAN DEFAULT FALSE,            -- Apakah menerapkan denda keterlambatan?
    late_fee_type VARCHAR(20) DEFAULT 'FIXED',     -- FIXED (Rupiah tetap) atau PERCENTAGE (Persentase)
    late_fee_value NUMERIC(12, 2) DEFAULT 0,       -- Nilai nominal denda (misal: Rp 10.000 atau 2.00%)
    grace_period_days INT DEFAULT 0,               -- Masa tenggang keterlambatan dalam hitungan hari
    
    -- Penunjang Aturan Kustom Dinamis Sekolah (JSONB)
    dynamic_rules JSONB DEFAULT '{}'::jsonb,       -- Menyimpan rule kustom seperti potongan bersyarat, dll.
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 6. TABEL PRODUK PEMBAYARAN (Payment Products)
CREATE TABLE payment_products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    policy_id UUID REFERENCES payment_policies(id) ON DELETE SET NULL, -- Merujuk ke template kebijakan pembayaran
    product_name VARCHAR(100) NOT NULL, -- Contoh: "SPP Kelas 10 - Juli 2026", "Uang Gedung Angkatan 2026"
    product_type VARCHAR(50) NOT NULL,  -- SPP, BUILDING_FEE, INFAQ, STUDY_TOUR, UNIFORM
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 7. TABEL KEWAJIBAN SISWA (Student Obligations - Inti Tagihan Murid)
CREATE TABLE student_obligations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES payment_products(id) ON DELETE RESTRICT,
    amount_total NUMERIC(12, 2) NOT NULL,     -- Nominal tagihan awal (misal: Rp 1.000.000)
    amount_remaining NUMERIC(12, 2) NOT NULL, -- Sisa tagihan aktif saat ini (berkurang seiring cicilan)
    due_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'UNPAID',       -- UNPAID, PARTIAL, PAID
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 8. TABEL NOTA TRANSAKSI UTAMA (Payment Orders - Invoice Header Checkout)
CREATE TABLE payment_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    payer_id UUID NOT NULL REFERENCES users(id), -- Orang tua yang checkout via Flutter ATAU Bendahara via Vue.js jika tunai
    total_amount NUMERIC(12, 2) NOT NULL,       -- Total keseluruhan dana yang dibayarkan dalam satu checkout
    payment_method VARCHAR(50) NOT NULL,        -- CASH, VA_BCA, VA_MANDIRI, GO_PAY, QRIS, dll.
    payment_status VARCHAR(20) DEFAULT 'PENDING',-- PENDING, SETTLED, EXPIRED, FAILED
    reference_number VARCHAR(100),              -- Nomor referensi/ID transaksi eksternal dari Payment Gateway atau nomor nota tunai
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 9. TABEL DETAIL ALOKASI PEMBAYARAN (Payment Allocations - Invoice Detail / Keranjang Belanja)
-- Menyimpan porsi nominal dana yang ditentukan secara manual (user-defined) untuk setiap item tagihan
CREATE TABLE payment_allocations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_order_id UUID NOT NULL REFERENCES payment_orders(id) ON DELETE CASCADE,
    obligation_id UUID NOT NULL REFERENCES student_obligations(id) ON DELETE RESTRICT,
    amount_allocated NUMERIC(12, 2) NOT NULL,  -- Nominal dana yang dialokasikan user ke kewajiban spesifik ini
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 10. TABEL BUKU BESAR KEUANGAN (Financial Ledgers - Histori Mutasi & Audit Trail)
-- Mencatat seluruh riwayat mutasi penambahan/pengurangan tagihan siswa secara detail dan mutlak (immutable)
CREATE TABLE financial_ledgers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    obligation_id UUID NOT NULL REFERENCES student_obligations(id) ON DELETE RESTRICT,
    payment_order_id UUID REFERENCES payment_orders(id) ON DELETE SET NULL, -- Kosong (NULL) jika berupa penyesuaian/debit awal tagihan
    transaction_type VARCHAR(10) NOT NULL,      -- DEBIT (Menambah tagihan siswa), CREDIT (Mengurangi tagihan siswa)
    amount NUMERIC(12, 2) NOT NULL,             -- Nominal transaksi mutasi
    balance_after NUMERIC(12, 2) NOT NULL,      -- Akumulasi sisa saldo kewajiban tepat setelah transaksi ini tercatat
    description TEXT,                           -- Contoh: "Pembayaran parsial via VA BCA", "Pencatatan tunai via bendahara", atau "Penyesuaian diskon beasiswa"
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);