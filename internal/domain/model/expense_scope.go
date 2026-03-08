package model

type ExpenseScope string

const (
	ExpenseScopeCommon           ExpenseScope = "Comú"
	ExpenseScopeOliveSection     ExpenseScope = "Secció d'oliva"
	ExpenseScopeLivestockSection ExpenseScope = "Secció de ramaderia"
	ExpenseScopePartner          ExpenseScope = "Soci"
)

func (es ExpenseScope) String() string {
	return string(es)
}
