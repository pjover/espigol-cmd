package model

import "testing"

func TestExpenseSubtypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		got      ExpenseSubtype
		expected string
	}{
		{"ExpenseSubtypeA1", ExpenseSubtypeA1, "[a1] Recerca, investigació, desenvolupament i innovació"},
		{"ExpenseSubtypeA2", ExpenseSubtypeA2, "[a2] Activitats d'informació i promoció de productes agraris"},
		{"ExpenseSubtypeA3", ExpenseSubtypeA3, "[a3] Activitats d'informació i promoció de productes agraris"},
		{"ExpenseSubtypeA4", ExpenseSubtypeA4, "[a4] Activitats de prevenció i control de malalties animals i plagues vegetals"},
		{"ExpenseSubtypeA5", ExpenseSubtypeA5, "[a5] Activitats que fomentin la sostenibilitat, el consum local i de quilòmetre zero"},
		{"ExpenseSubtypeA6", ExpenseSubtypeA6, "[a6] Despeses de fertilitzants, productes d'alimentació animal i ormejos"},
		{"ExpenseSubtypeB1", ExpenseSubtypeB1, "[b1] Despeses d'adquisició de maquinària i materials"},
		{"ExpenseSubtypeB2", ExpenseSubtypeB2, "[b2] Manteniment i restauració d'elements etnològics vinculats al manteniment de l'activitat agrària"},
		{"ExpenseSubtypeB3", ExpenseSubtypeB3, "[b3] Inversions en elements, immobles o estructures permanents pel desplegament, promoció i venda dels productes agraris"},
		{"ExpenseSubtypeB4", ExpenseSubtypeB4, "[b4] Despeses de registre de posicionament de productes"},
		{"ExpenseSubtypeB5", ExpenseSubtypeB5, "[b5] Accions de millora de la imatge dels productes propis"},
		{"ExpenseSubtypeC1", ExpenseSubtypeC1, "[c1] Despeses estructurals"},
		{"ExpenseSubtypeC2", ExpenseSubtypeC2, "[c2] Despeses de personal"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.got) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.got))
			}
		})
	}
}

func TestExpenseSubtypeType(t *testing.T) {
	tests := []struct {
		subtype    ExpenseSubtype
		parentType ExpenseType
	}{
		{ExpenseSubtypeA1, ExpenseTypeA},
		{ExpenseSubtypeA2, ExpenseTypeA},
		{ExpenseSubtypeA3, ExpenseTypeA},
		{ExpenseSubtypeA4, ExpenseTypeA},
		{ExpenseSubtypeA5, ExpenseTypeA},
		{ExpenseSubtypeA6, ExpenseTypeA},
		{ExpenseSubtypeB1, ExpenseTypeB},
		{ExpenseSubtypeB2, ExpenseTypeB},
		{ExpenseSubtypeB3, ExpenseTypeB},
		{ExpenseSubtypeB4, ExpenseTypeB},
		{ExpenseSubtypeB5, ExpenseTypeB},
		{ExpenseSubtypeC1, ExpenseTypeC},
		{ExpenseSubtypeC2, ExpenseTypeC},
	}
	
	for _, tt := range tests {
		t.Run(string(tt.subtype), func(t *testing.T) {
			got := tt.subtype.Type()
			if got != tt.parentType {
				t.Errorf("Expected parent type %s, got %s", tt.parentType, got)
			}
		})
	}
}

func TestExpenseSubtypeString(t *testing.T) {
	tests := []struct {
		subtype  ExpenseSubtype
		expected string
	}{
		{ExpenseSubtypeA1, "[a1] Recerca, investigació, desenvolupament i innovació"},
		{ExpenseSubtypeB2, "[b2] Manteniment i restauració d'elements etnològics vinculats al manteniment de l'activitat agrària"},
		{ExpenseSubtypeC2, "[c2] Despeses de personal"},
	}
	
	for _, tt := range tests {
		t.Run(string(tt.subtype), func(t *testing.T) {
			got := tt.subtype.String()
			if got != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, got)
			}
		})
	}
}
