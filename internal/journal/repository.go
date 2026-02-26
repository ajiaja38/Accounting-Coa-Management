package journal

import (
	"fiber.com/session-api/internal/domain"
	"fiber.com/session-api/pkg/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JournalDetailRow struct {
	ID          string  `gorm:"column:id"`
	CoaCode     string  `gorm:"column:coa_code"`
	CoaName     string  `gorm:"column:coa_name"`
	Debit       float64 `gorm:"column:debit"`
	Credit      float64 `gorm:"column:credit"`
	Description string  `gorm:"column:description"`
}

type Repository interface {
	FindAll(req *model.PaginationRequest) ([]JournalListResponse, int64, error)
	FindByID(id uuid.UUID) (*domain.JournalEntry, []JournalDetailRow, error)
	Create(entry *domain.JournalEntry, details []domain.JournalEntryDetail) error
	PostJournal(id uuid.UUID) error
	Delete(id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(req *model.PaginationRequest) ([]JournalListResponse, int64, error) {
	var total int64
	offset := (req.Page - 1) * req.Limit
	search := "%" + req.Search + "%"

	jeWhere := "je.deleted_at IS NULL"

	if req.Search != "" {
		jeWhere += " AND (je.reference ILIKE ? OR je.description ILIKE ?)"
	}

	countQuery := `
		SELECT COUNT(*)
		FROM journal_entries je
		WHERE ` + jeWhere

	var args []interface{}
	if req.Search != "" {
		args = append(args, search, search)
	}

	if err := r.db.Raw(countQuery, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	dataQuery := `
		SELECT
			je.id,
			je.date,
			je.reference,
			je.description,
			je.status,
			je.created_by,
			COALESCE(SUM(jd.debit), 0)  AS total_debit,
			COALESCE(SUM(jd.credit), 0) AS total_credit
		FROM journal_entries je
		LEFT JOIN journal_entry_details jd
			ON jd.journal_entry_id = je.id
			AND jd.deleted_at IS NULL
		WHERE ` + jeWhere + `
		GROUP BY
			je.id,
			je.date,
			je.reference,
			je.description,
			je.status,
			je.created_by
		ORDER BY
			je.date DESC,
			je.created_at DESC
		LIMIT ? OFFSET ?`

	argsWithLimit := append(args, req.Limit, offset)

	var rows []JournalListResponse
	if err := r.db.Raw(dataQuery, argsWithLimit...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *repository) FindByID(id uuid.UUID) (*domain.JournalEntry, []JournalDetailRow, error) {
	var entry domain.JournalEntry
	result := r.db.Raw(
		`SELECT
			id,
			date,
			reference,
			description,
			status,
			created_by,
			created_at,
			updated_at
		 FROM journal_entries
		 WHERE id = ?
		 AND deleted_at IS NULL
		 LIMIT 1`,
		id,
	).Scan(&entry)

	if result.Error != nil {
		return nil, nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil, nil
	}

	var details []JournalDetailRow
	detailQuery := `
		SELECT
			jd.id,
			jd.coa_code,
			c.name AS coa_name,
			jd.debit,
			jd.credit,
			jd.description
		FROM journal_entry_details jd
		JOIN chart_of_accounts c ON c.code = jd.coa_code
		WHERE jd.journal_entry_id = ?
		AND jd.deleted_at IS NULL
		ORDER BY jd.debit DESC
	`

	if err := r.db.Raw(detailQuery, id).Scan(&details).Error; err != nil {
		return nil, nil, err
	}

	return &entry, details, nil
}

func (r *repository) Create(entry *domain.JournalEntry, details []domain.JournalEntryDetail) error {
	if err := r.db.Exec(
		`INSERT INTO journal_entries (id, date, reference, description, status, created_by, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`,
		entry.ID, entry.Date, entry.Reference, entry.Description, entry.Status, entry.CreatedBy,
	).Error; err != nil {
		return err
	}

	for _, detail := range details {
		if err := r.db.Exec(
			`INSERT INTO journal_entry_details (id, journal_entry_id, coa_code, debit, credit, description)
			 VALUES (gen_random_uuid(), ?, ?, ?, ?, ?)`,
			detail.JournalEntryID, detail.CoaCode, detail.Debit, detail.Credit, detail.Description,
		).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) PostJournal(id uuid.UUID) error {
	result := r.db.Exec(
		`UPDATE journal_entries SET status = 'posted', updated_at = NOW()
		 WHERE id = ? AND status = 'draft' AND deleted_at IS NULL`,
		id,
	)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Journal not found or already posted")
	}
	return nil
}

func (r *repository) Delete(id uuid.UUID) error {
	result := r.db.Exec(
		`UPDATE journal_entries SET deleted_at = NOW()
		 WHERE id = ? AND status = 'draft' AND deleted_at IS NULL`,
		id,
	)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Journal not found or cannot be deleted (only draft journals can be deleted)")
	}
	return nil
}
