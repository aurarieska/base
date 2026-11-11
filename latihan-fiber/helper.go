package main

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"latihan-fiber/app/model"
)

// reqCtx memberi batas waktu untuk setiap operasi basis data.
// Tanpa batas waktu, satu query yang menggantung dapat menahan koneksi
// selamanya dan lama-lama menghabiskan seluruh isi pool.
func reqCtx(c *fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.UserContext(), 5*time.Second)
}

// paramID membaca path parameter :id dan memastikan berupa angka positif.
func paramID(c *fiber.Ctx) (int, bool) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

// parseListQuery membaca seluruh parameter query untuk endpoint daftar.
// Setiap parameter memiliki nilai bawaan yang masuk akal: permintaan
// daftar tidak boleh gagal hanya karena klien lupa mengirim limit.
func parseListQuery(c *fiber.Ctx) model.ListQuery {
	q := model.ListQuery{
		Search: strings.TrimSpace(c.Query("search")),
		Sort:   strings.ToLower(c.Query("sort", "id")),
		Order:  strings.ToLower(c.Query("order", "asc")),
		Page:   atoiOr(c.Query("page"), 1),
		Limit:  atoiOr(c.Query("limit"), 10),
	}

	// Pemeriksaan kasar di sini; daftar putih sebenarnya berada di
	// repository. Bila klien mengirim nilai di luar yang dikenal,
	// jatuhkan ke "id" sebagai bawaan.
	switch q.Sort {
	case "id", "nim", "name", "created_at":
		// biarkan
	default:
		q.Sort = "id"
	}
	if q.Order != "desc" {
		q.Order = "asc"
	}
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 10
	}

	// Filter is_active hanya aktif bila klien mengirim "true" atau "false".
	if raw := c.Query("is_active"); raw != "" {
		v := strings.EqualFold(raw, "true")
		q.IsActive = &v
	}

	return q
}

func atoiOr(s string, fallback int) int {
	if s == "" {
		return fallback
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return n
}

// requireJSON memastikan content-type pada permintaan ber-body adalah JSON.
// Dipasang sebagai middleware di grup /students.
func requireJSON(c *fiber.Ctx) error {
	method := c.Method()
	if method == fiber.MethodPost || method == fiber.MethodPut || method == fiber.MethodPatch {
		ct := strings.ToLower(c.Get(fiber.HeaderContentType))
		if !strings.Contains(ct, "application/json") {
			return fail(c, fiber.StatusUnsupportedMediaType,
				"content-type harus application/json")
		}
	}
	return c.Next()
}

// Pembalas respons: satu tempat untuk membangun bentuk WebResponse,
// sehingga handler bisa fokus pada logika bisnisnya saja.

func ok(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(model.WebResponse{
		Message: message,
		Data:    data,
	})
}

func okList(c *fiber.Ctx, message string, data interface{}, meta *model.Meta) error {
	return c.Status(fiber.StatusOK).JSON(model.WebResponse{
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func created(c *fiber.Ctx, message string, data interface{}, location string) error {
	if location != "" {
		c.Set("Location", location)
	}
	return c.Status(fiber.StatusCreated).JSON(model.WebResponse{
		Message: message,
		Data:    data,
	})
}

func noContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func fail(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(model.WebResponse{
		Message: message,
	})
}

func failValidation(c *fiber.Ctx, errs map[string]string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(model.WebResponse{
		Message: "validasi gagal",
		Errors:  errs,
	})
}