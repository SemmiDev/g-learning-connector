package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
	gl "lab.garudacyber.co.id/g-learning-connector"
	"net/http"
)

type (
	GetTotalStudentsResponse struct {
		Total int64 `json:"total"`
	}

	ListStudentsRequest struct {
		gl.Filter
	}

	ListStudentsResponse struct {
		ID        string `json:"id" gorm:"column:id"`
		Name      string `json:"name" gorm:"column:nama_mahasiswa"`
		Gender    string `json:"gender" gorm:"column:jenis_kelamin"`
		NIK       string `json:"nik" gorm:"column:nik"`
		Email     string `json:"email" gorm:"column:email"`
		Handphone string `json:"handphone" gorm:"column:handphone"`
		Telephone string `json:"telephone" gorm:"column:telepon"`
	}
)

func NewListStudentsRequest() *ListStudentsRequest {
	return &ListStudentsRequest{
		Filter: gl.NewFilterPagination(),
	}
}

func (a *ApplicationServer) ListStudents(c *fiber.Ctx) error {
	req := NewListStudentsRequest()
	if err := c.QueryParser(req); err != nil {
		return HandleError(c, err)
	}

	listStudents := make([]ListStudentsResponse, 0)

	offset := req.Filter.GetOffset()
	limit := req.Filter.GetLimit()

	q := a.db.
		Select(`id, nama_mahasiswa, jenis_kelamin, nik, email, handphone, telepon`).
		Table("mahasiswa")

	if req.Filter.HasKeyword() {
		q = q.Where("nama_mahasiswa LIKE ? OR nik LIKE ?", "%"+req.Filter.Keyword+"%", "%"+req.Filter.Keyword+"%")
	}

	if req.Filter.HasSort() {
		q = q.Order(
			clause.OrderByColumn{
				Column: clause.Column{Name: req.Filter.SortBy},
				Desc:   req.Filter.IsDesc(),
			},
		)
	} else {
		q = q.Order("created_at ASC")
	}

	// Menghitung jumlah total data tanpa offset dan limit
	var totalData int64
	if err := q.Count(&totalData).Error; err != nil {
		return HandleError(c, err)
	}

	// Menambahkan limit dan offset setelah menghitung total data
	q = q.Offset(int(offset)).Limit(int(limit))

	if err := q.Scan(&listStudents).Error; err != nil {
		return HandleError(c, err)
	}

	pageInfo, err := gl.NewPageInfo(req.Filter.CurrentPage, limit, offset, totalData)
	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(ApiResponse[ListDataApiResponseWrapper[ListStudentsResponse]]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan data mahasiswa",
		Data: ListDataApiResponseWrapper[ListStudentsResponse]{
			List:     listStudents,
			PageInfo: pageInfo,
		},
	})
}

func (a *ApplicationServer) GetTotalStudents(c *fiber.Ctx) error {
	var total int64

	err := a.db.
		Table("mahasiswa").
		Where("mahasiswa.deleted_at IS NULL").
		Count(&total).
		Error

	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(ApiResponse[GetTotalStudentsResponse]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan total mahasiswa",
		Data: GetTotalStudentsResponse{
			Total: total,
		},
	})
}
