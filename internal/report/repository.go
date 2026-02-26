package report

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetOpeningBalance(coaCode, startDate string) (float64, float64, error)
	GetLedgerTransactions(coaCode, startDate, endDate string) ([]TransactionRow, error)
	GetAccountBalances(startDate, endDate string) ([]AccountBalanceRow, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetOpeningBalance(coaCode, startDate string) (float64, float64, error) {
	if startDate == "" {
		return 0, 0, nil
	}
	var res struct {
		Debit  float64
		Credit float64
	}
	query := `
		SELECT 
			COALESCE(SUM(jd.debit), 0) as debit, 
			COALESCE(SUM(jd.credit), 0) as credit
		FROM journal_entry_details jd
		JOIN journal_entries je ON je.id = jd.journal_entry_id
		WHERE jd.coa_code = ? 
		  AND je.status = 'posted' 
		  AND je.deleted_at IS NULL
		  AND je.date < ?
	`
	if err := r.db.Raw(query, coaCode, startDate).Scan(&res).Error; err != nil {
		return 0, 0, err
	}
	return res.Debit, res.Credit, nil
}

func (r *repository) GetLedgerTransactions(coaCode, startDate, endDate string) ([]TransactionRow, error) {
	var rows []TransactionRow

	query := `
		SELECT 
			je.date, 
			je.reference, 
			jd.description, 
			jd.debit, 
			jd.credit
		FROM journal_entry_details jd
		JOIN journal_entries je ON je.id = jd.journal_entry_id
		WHERE jd.coa_code = ? 
		  AND je.status = 'posted' 
		  AND je.deleted_at IS NULL
	`
	var args []interface{}
	args = append(args, coaCode)

	if startDate != "" {
		query += " AND je.date >= ?"
		args = append(args, startDate)
	}
	if endDate != "" {
		query += " AND je.date <= ?"
		args = append(args, endDate)
	}
	query += " ORDER BY je.date ASC, je.created_at ASC"

	if err := r.db.Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *repository) GetAccountBalances(startDate, endDate string) ([]AccountBalanceRow, error) {
	var rows []AccountBalanceRow

	onClause := "je.id = jd.journal_entry_id AND je.status = 'posted' AND je.deleted_at IS NULL"
	var args []any

	if startDate != "" {
		onClause += " AND je.date >= ?"
		args = append(args, startDate)
	}
	if endDate != "" {
		onClause += " AND je.date <= ?"
		args = append(args, endDate)
	}

	query := `
		SELECT 
			c.code AS coa_code,
			c.name AS coa_name,
			c.type AS type,
			COALESCE(SUM(jd.debit), 0) AS sum_debit,
			COALESCE(SUM(jd.credit), 0) AS sum_credit
		FROM chart_of_accounts c
		LEFT JOIN journal_entry_details jd ON jd.coa_code = c.code AND jd.deleted_at IS NULL
		LEFT JOIN journal_entries je ON ` + onClause + `
		WHERE c.is_active = true
		GROUP BY c.code, c.name, c.type
		ORDER BY c.code ASC
	`

	if err := r.db.Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
