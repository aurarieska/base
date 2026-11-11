-- Migrasi pertama: membuat tabel students.
-- Satu baris mewakili satu mahasiswa yang dicatat di sistem.
CREATE TABLE IF NOT EXISTS students (
    id         SERIAL       PRIMARY KEY,
    nim        VARCHAR(10)  NOT NULL,
    name       VARCHAR(100) NOT NULL,
    grade      INTEGER      NOT NULL DEFAULT 0,
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    -- Nilai students.grade berupa angka 0 sampai 100.
    -- Pemeriksaan ini dijaga basis data agar tidak ada jalan masuk lain
    -- yang bisa menitipkan nilai di luar rentang itu.
    CONSTRAINT students_grade_range CHECK (grade BETWEEN 0 AND 100)
);

-- Keunikan NIM tanpa membedakan huruf besar dan kecil.
-- Dibuat sebagai indeks pada ekspresi LOWER(nim) supaya pemeriksaan
-- pencocokan "12345" dengan "12345" tetap konsisten bila kelak kolom
-- diperbesar dan menerima huruf.
CREATE UNIQUE INDEX IF NOT EXISTS students_nim_lower_key
    ON students (LOWER(nim));

-- Indeks tambahan pada name untuk mempercepat pencarian ILIKE.
-- Tanpa indeks, pencarian pada tabel besar akan membaca seluruh baris.
CREATE INDEX IF NOT EXISTS students_name_lower_idx
    ON students (LOWER(name));