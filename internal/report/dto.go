package report

import "time"

// LedgerQuery is the request DTO for General Ledger.
// CoaCode is required. StartDate and EndDate are optional.
type LedgerQuery struct {
	CoaCode   string `query:"coaCode" validate:"required"`
	StartDate string `query:"startDate"`
	EndDate   string `query:"endDate"`
}

// PeriodQuery is the request DTO for periodic reports.
type PeriodQuery struct {
	StartDate string `query:"startDate"`
	EndDate   string `query:"endDate"`
}

// TransactionRow represents a single line in the general ledger.
type TransactionRow struct {
	Date        time.Time `json:"date"`
	Reference   string    `json:"reference"`
	Description string    `json:"description"`
	Debit       float64   `json:"debit"`
	Credit      float64   `json:"credit"`
	Balance     float64   `json:"balance"` // calculated running balance
}

// LedgerResponse is the response body for General Ledger.
type LedgerResponse struct {
	CoaCode        string           `json:"coaCode"`
	CoaName        string           `json:"coaName"`
	OpeningBalance float64          `json:"openingBalance"`
	Transactions   []TransactionRow `json:"transactions"`
	ClosingBalance float64          `json:"closingBalance"`
}

// AccountBalanceRow represents a summarized account balance for a period.
type AccountBalanceRow struct {
	CoaCode string  `json:"coaCode" gorm:"column:coa_code"`
	CoaName string  `json:"coaName" gorm:"column:coa_name"`
	Type    string  `json:"-"       gorm:"column:type"`
	Debit   float64 `json:"debit,omitempty"   gorm:"column:sum_debit"`
	Credit  float64 `json:"credit,omitempty"  gorm:"column:sum_credit"`
	Balance float64 `json:"balance,omitempty"` // Net balance for PnL/BalanceSheet
}

// TrialBalanceResponse is the response body for Trial Balance.
type TrialBalanceResponse struct {
	Rows        []AccountBalanceRow `json:"rows"`
	TotalDebit  float64             `json:"totalDebit"`
	TotalCredit float64             `json:"totalCredit"`
	IsBalanced  bool                `json:"isBalanced"`
}

// ProfitLossResponse is the response body for PnL.
type ProfitLossResponse struct {
	Revenues     []AccountBalanceRow `json:"revenues"`
	TotalRevenue float64             `json:"totalRevenue"`
	Expenses     []AccountBalanceRow `json:"expenses"`
	TotalExpense float64             `json:"totalExpense"`
	NetProfit    float64             `json:"netProfit"`
}

// BalanceSheetResponse is the response body for Balance Sheet.
type BalanceSheetResponse struct {
	Assets          []AccountBalanceRow `json:"assets"`
	TotalAsset      float64             `json:"totalAsset"`
	Liabilities     []AccountBalanceRow `json:"liabilities"`
	TotalLiability  float64             `json:"totalLiability"`
	Equities        []AccountBalanceRow `json:"equities"`
	TotalEquity     float64             `json:"totalEquity"`
	TotalLiabEquity float64             `json:"totalLiabilityAndEquity"`
	IsBalanced      bool                `json:"isBalanced"`
}

// Swagger Responses

type SwaggerLedgerResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    LedgerResponse `json:"data"`
}

type SwaggerTrialBalanceResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    TrialBalanceResponse `json:"data"`
}

type SwaggerProfitLossResponse struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    ProfitLossResponse `json:"data"`
}

type SwaggerBalanceSheetResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    BalanceSheetResponse `json:"data"`
}
