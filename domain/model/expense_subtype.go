package model

type ExpenseSubtype string

const (
	// Type A: Current expenses
	ExpenseSubtypeA1 ExpenseSubtype = "[a1] Recerca, investigació, desenvolupament i innovació"
	ExpenseSubtypeA2 ExpenseSubtype = "[a2] Activitats d'informació i promoció de productes agraris"
	ExpenseSubtypeA3 ExpenseSubtype = "[a3] Activitats d'informació i promoció de productes agraris"
	ExpenseSubtypeA4 ExpenseSubtype = "[a4] Activitats de prevenció i control de malalties animals i plagues vegetals"
	ExpenseSubtypeA5 ExpenseSubtype = "[a5] Activitats que fomentin la sostenibilitat, el consum local i de quilòmetre zero"
	ExpenseSubtypeA6 ExpenseSubtype = "[a6] Despeses de fertilitzants, productes d'alimentació animal i ormejos"

	// Type B: Permanent or investment expenses
	ExpenseSubtypeB1 ExpenseSubtype = "[b1] Despeses d'adquisició de maquinària i materials"
	ExpenseSubtypeB2 ExpenseSubtype = "[b2] Manteniment i restauració d'elements etnològics vinculats al manteniment de l'activitat agrària"
	ExpenseSubtypeB3 ExpenseSubtype = "[b3] Inversions en elements, immobles o estructures permanents pel desplegament, promoció i venda dels productes agraris"
	ExpenseSubtypeB4 ExpenseSubtype = "[b4] Despeses de registre de posicionament de productes"
	ExpenseSubtypeB5 ExpenseSubtype = "[b5] Accions de millora de la imatge dels productes propis"

	// Type C: Indirect or structural expenses
	ExpenseSubtypeC1 ExpenseSubtype = "[c1] Despeses estructurals"
	ExpenseSubtypeC2 ExpenseSubtype = "[c2] Despeses de personal"
)

func (es ExpenseSubtype) Type() ExpenseType {
	switch es {
	case ExpenseSubtypeA1, ExpenseSubtypeA2, ExpenseSubtypeA3, ExpenseSubtypeA4, ExpenseSubtypeA5, ExpenseSubtypeA6:
		return ExpenseTypeA
	case ExpenseSubtypeB1, ExpenseSubtypeB2, ExpenseSubtypeB3, ExpenseSubtypeB4, ExpenseSubtypeB5:
		return ExpenseTypeB
	case ExpenseSubtypeC1, ExpenseSubtypeC2:
		return ExpenseTypeC
	default:
		return ""
	}
}

func (es ExpenseSubtype) String() string {
	return string(es)
}
