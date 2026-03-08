package model

type ExpenseCategory string

const (
	ExpenseCategoryCurrent    ExpenseCategory = "Despesa corrent"
	ExpenseCategoryInvestment ExpenseCategory = "Despesa d'inversió"
)

func (ec ExpenseCategory) String() string {
	return string(ec)
}

func (et ExpenseType) Category() ExpenseCategory {
	switch et {
	case ExpenseTypeA, ExpenseTypeC:
		return ExpenseCategoryCurrent
	case ExpenseTypeB:
		return ExpenseCategoryInvestment
	default:
		return ""
	}
}
