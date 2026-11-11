package model

import "time"

// Student adalah bentuk data student seperti yang tersimpan di tabel.
// Tag json menentukan nama field pada respons HTTP.
type Student struct {
	ID        int       `json:"id"`
	NIM       string    `json:"nim"`
	Name      string    `json:"name"`
	Grade     int       `json:"grade"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateStudentRequest adalah bentuk masukan untuk POST /students.
// Field yang wajib: nim, name. Grade boleh 0 (bawaan), is_active selalu true.
type CreateStudentRequest struct {
	NIM   string `json:"nim"`
	Name  string `json:"name"`
	Grade int    `json:"grade"`
}

// ReplaceStudentRequest adalah bentuk masukan untuk PUT /students/:id.
// Seluruh field wajib diisi, karena PUT mengganti seluruh baris.
type ReplaceStudentRequest struct {
	NIM      string `json:"nim"`
	Name     string `json:"name"`
	Grade    int    `json:"grade"`
	IsActive bool   `json:"is_active"`
}

// PatchStudentRequest adalah bentuk masukan untuk PATCH /students/:id.
// Field berupa pointer agar "tidak ada di JSON" dapat dibedakan dari
// "ada tetapi nilainya nol". Hanya field yang dikirim klien yang diubah.
type PatchStudentRequest struct {
	NIM      *string `json:"nim,omitempty"`
	Name     *string `json:"name,omitempty"`
	Grade    *int    `json:"grade,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// ListQuery memuat seluruh parameter kueri untuk GET /students.
// Dipakai bersama antara handler dan repository agar logika pembangunan
// query terpusat di satu tempat (di repository).
type ListQuery struct {
	Search   string // dicari dengan ILIKE pada kolom name
	IsActive *bool  // nil berarti tidak memfilter
	Sort     string // nama kolom dari daftar putih
	Order    string // "asc" atau "desc"
	Page     int    // dimulai dari 1
	Limit    int    // jumlah baris per halaman
}

// Offset menghitung berapa baris yang dilewati untuk halaman ini.
// Perhitungan ini pindah ke sini karena kini dipakai langsung oleh SQL.
func (q ListQuery) Offset() int {
	return (q.Page - 1) * q.Limit
}

// Meta memuat informasi paginasi pada respons daftar.
// total berasal dari SELECT COUNT(*) dengan syarat penyaringan yang sama.
type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// WebResponse adalah bentuk baku seluruh respons API.
// Field opsional memakai omitempty agar tidak muncul bila bernilai nol.
type WebResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}