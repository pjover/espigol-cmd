package model

type ExpenseType string

const (
	ExpenseTypeA ExpenseType = "[a] Despeses corrents de caràcter fungible o temporals"
	ExpenseTypeB ExpenseType = "[b] Despeses de caràcter permanent o d'inversió"
	ExpenseTypeC ExpenseType = "[c] Despeses indirectes o d'estructura de l'entitat beneficiària"
)

func (et ExpenseType) String() string {
	return string(et)
}
