package coa

import (
	"fiber.com/session-api/internal/domain"
	"fiber.com/session-api/pkg/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll(req *model.PaginationRequest) ([]domain.ChartOfAccount, int64, error)
	FindAllWithChildren(req *model.PaginationRequest) ([]CoaReqursiveResponse, int64, error)
	FindByCode(code string) (*domain.ChartOfAccount, error)
	Create(coa *domain.ChartOfAccount) error
	Update(coa *domain.ChartOfAccount) error
	Delete(code string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(req *model.PaginationRequest) ([]domain.ChartOfAccount, int64, error) {
	var accounts []domain.ChartOfAccount
	var total int64
	offset := (req.Page - 1) * req.Limit
	search := "%" + req.Search + "%"

	countQuery := `SELECT COUNT(*) FROM chart_of_accounts WHERE deleted_at IS NULL AND (name ILIKE ? OR code ILIKE ?)`
	if err := r.db.Raw(countQuery, search, search).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	dataQuery := `
		SELECT code, name, type, parent_code, is_active, created_at, updated_at
		FROM chart_of_accounts
		WHERE deleted_at IS NULL AND (name ILIKE ? OR code ILIKE ?)
		ORDER BY code ASC
		LIMIT ? OFFSET ?`

	if err := r.db.Raw(dataQuery, search, search, req.Limit, offset).Scan(&accounts).Error; err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

func (r *repository) FindAllWithChildren(req *model.PaginationRequest) ([]CoaReqursiveResponse, int64, error) {
	var accounts []CoaReqursiveResponse
	var total int64
	offset := (req.Page - 1) * req.Limit
	search := "%" + req.Search + "%"
	var query string
	var args []any

	if req.Search != "" {
		query = `
            and (
                coa."name" ilike ?
                OR coa."type" ilike ?
            )
        `
		args = append(args, search, search)
	}

	countQuery := `
        SELECT COUNT(*)
        from chart_of_accounts coa
        WHERE coa.parent_code IS NULL
        AND coa.is_active = true
        ` + query + `
    `

	if err := r.db.Raw(countQuery, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	dataQuery := `
        SELECT
            coa.code,
            coa."name",
            coa."type",
            COALESCE(
                JSON_AGG(
                    JSON_BUILD_OBJECT(
                        'code', coa2.code,
                        'name', coa2."name",
                        'type', coa2."type",
                        'parentCode', coa2.parent_code
                    )
                ) FILTER (WHERE coa2.code is not null),
                '[]'::json
            ) as childs
        FROM chart_of_accounts coa
        LEFT JOIN chart_of_accounts coa2
            on coa2.parent_code = coa.code
            and coa2.is_active = true
        WHERE coa.parent_code is null
        AND coa.is_active = true
        ` + query + `
        GROUP BY
            coa.code,
            coa."name",
            coa."type"
        ORDER BY coa.created_at asc
        LIMIT ?
        OFFSET ?
	`
	args = append(args, req.Limit, offset)

	if err := r.db.Raw(dataQuery, args...).Scan(&accounts).Error; err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

func (r *repository) FindByCode(code string) (*domain.ChartOfAccount, error) {
	var coa domain.ChartOfAccount
	result := r.db.Raw(
		`SELECT code, name, type, parent_code, is_active, created_at, updated_at
		 FROM chart_of_accounts WHERE code = ? AND deleted_at IS NULL LIMIT 1`,
		code,
	).Scan(&coa)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &coa, nil
}

func (r *repository) Create(coa *domain.ChartOfAccount) error {
	return r.db.Exec(
		`INSERT INTO chart_of_accounts (code, name, type, parent_code, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, NOW(), NOW())`,
		coa.Code, coa.Name, coa.Type, coa.ParentCode, coa.IsActive,
	).Error
}

func (r *repository) Update(coa *domain.ChartOfAccount) error {
	return r.db.Exec(
		`UPDATE chart_of_accounts
		 SET name = ?, type = ?, parent_code = ?, is_active = ?, updated_at = NOW()
		 WHERE code = ? AND deleted_at IS NULL`,
		coa.Name, coa.Type, coa.ParentCode, coa.IsActive, coa.Code,
	).Error
}

func (r *repository) Delete(code string) error {
	result := r.db.Exec(
		`UPDATE chart_of_accounts SET deleted_at = NOW() WHERE code = ? AND deleted_at IS NULL`,
		code,
	)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "COA not found")
	}
	return nil
}
