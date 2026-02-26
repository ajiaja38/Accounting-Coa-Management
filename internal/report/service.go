package report

import (
	"fiber.com/session-api/internal/coa"
	"fiber.com/session-api/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type Service interface {
	GetLedger(req *LedgerQuery) (*LedgerResponse, error)
	GetTrialBalance(req *PeriodQuery) (*TrialBalanceResponse, error)
	GetProfitLoss(req *PeriodQuery) (*ProfitLossResponse, error)
	GetBalanceSheet(req *PeriodQuery) (*BalanceSheetResponse, error)
}

type service struct {
	repo    Repository
	coaRepo coa.Repository
}

func NewService(repo Repository, coaRepo coa.Repository) Service {
	return &service{repo: repo, coaRepo: coaRepo}
}

func (s *service) GetLedger(req *LedgerQuery) (*LedgerResponse, error) {
	account, err := s.coaRepo.FindByCode(req.CoaCode)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if account == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Account not found")
	}

	debit, credit, err := s.repo.GetOpeningBalance(req.CoaCode, req.StartDate)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	openingBalance := 0.0
	if account.Type == domain.AccountTypeAsset || account.Type == domain.AccountTypeExpense {
		openingBalance = debit - credit
	} else {
		openingBalance = credit - debit
	}

	transactions, err := s.repo.GetLedgerTransactions(req.CoaCode, req.StartDate, req.EndDate)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if transactions == nil {
		transactions = []TransactionRow{}
	}

	currentBalance := openingBalance
	for i := range transactions {
		if account.Type == domain.AccountTypeAsset || account.Type == domain.AccountTypeExpense {
			currentBalance += transactions[i].Debit - transactions[i].Credit
		} else {
			currentBalance += transactions[i].Credit - transactions[i].Debit
		}
		transactions[i].Balance = currentBalance
	}

	return &LedgerResponse{
		CoaCode:        account.Code,
		CoaName:        account.Name,
		OpeningBalance: openingBalance,
		Transactions:   transactions,
		ClosingBalance: currentBalance,
	}, nil
}

func (s *service) GetTrialBalance(req *PeriodQuery) (*TrialBalanceResponse, error) {
	balances, err := s.repo.GetAccountBalances(req.StartDate, req.EndDate)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if balances == nil {
		balances = []AccountBalanceRow{}
	}

	totalDebit := 0.0
	totalCredit := 0.0

	for _, bal := range balances {
		totalDebit += bal.Debit
		totalCredit += bal.Credit
	}

	isBalanced := totalDebit == totalCredit

	return &TrialBalanceResponse{
		Rows:        balances,
		TotalDebit:  totalDebit,
		TotalCredit: totalCredit,
		IsBalanced:  isBalanced,
	}, nil
}

func (s *service) GetProfitLoss(req *PeriodQuery) (*ProfitLossResponse, error) {
	balances, err := s.repo.GetAccountBalances(req.StartDate, req.EndDate)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	res := &ProfitLossResponse{
		Revenues: []AccountBalanceRow{},
		Expenses: []AccountBalanceRow{},
	}

	for _, bal := range balances {
		if bal.Type == string(domain.AccountTypeRevenue) {
			net := bal.Credit - bal.Debit
			if net != 0 {
				res.Revenues = append(res.Revenues, AccountBalanceRow{
					CoaCode: bal.CoaCode,
					CoaName: bal.CoaName,
					Balance: net,
				})
				res.TotalRevenue += net
			}
		} else if bal.Type == string(domain.AccountTypeExpense) {
			net := bal.Debit - bal.Credit
			if net != 0 {
				res.Expenses = append(res.Expenses, AccountBalanceRow{
					CoaCode: bal.CoaCode,
					CoaName: bal.CoaName,
					Balance: net,
				})
				res.TotalExpense += net
			}
		}
	}

	res.NetProfit = res.TotalRevenue - res.TotalExpense
	return res, nil
}

func (s *service) GetBalanceSheet(req *PeriodQuery) (*BalanceSheetResponse, error) {
	balances, err := s.repo.GetAccountBalances("", req.EndDate)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	totalRev := 0.0
	totalExp := 0.0

	res := &BalanceSheetResponse{
		Assets:      []AccountBalanceRow{},
		Liabilities: []AccountBalanceRow{},
		Equities:    []AccountBalanceRow{},
	}

	for _, bal := range balances {
		if bal.Type == string(domain.AccountTypeAsset) {
			net := bal.Debit - bal.Credit
			if net != 0 {
				res.Assets = append(res.Assets, AccountBalanceRow{
					CoaCode: bal.CoaCode,
					CoaName: bal.CoaName,
					Balance: net,
				})
				res.TotalAsset += net
			}
		} else if bal.Type == string(domain.AccountTypeLiability) {
			net := bal.Credit - bal.Debit
			if net != 0 {
				res.Liabilities = append(res.Liabilities, AccountBalanceRow{
					CoaCode: bal.CoaCode,
					CoaName: bal.CoaName,
					Balance: net,
				})
				res.TotalLiability += net
			}
		} else if bal.Type == string(domain.AccountTypeEquity) {
			net := bal.Credit - bal.Debit
			if net != 0 {
				res.Equities = append(res.Equities, AccountBalanceRow{
					CoaCode: bal.CoaCode,
					CoaName: bal.CoaName,
					Balance: net,
				})
				res.TotalEquity += net
			}
		} else if bal.Type == string(domain.AccountTypeRevenue) {
			totalRev += bal.Credit - bal.Debit
		} else if bal.Type == string(domain.AccountTypeExpense) {
			totalExp += bal.Debit - bal.Credit
		}
	}

	netProfit := totalRev - totalExp
	res.Equities = append(res.Equities, AccountBalanceRow{
		CoaCode: "-",
		CoaName: "Laba Periode Berjalan (Net Profit)",
		Balance: netProfit,
	})
	res.TotalEquity += netProfit

	res.TotalLiabEquity = res.TotalLiability + res.TotalEquity
	res.IsBalanced = res.TotalAsset == res.TotalLiabEquity

	return res, nil
}
