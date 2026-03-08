package cfg

var defaultValues = map[string]interface{}{
	"db.server":                       "mongodb://localhost:27017",
	"db.name":                         "espigol",
	"server.port":                     "8080",
	"expenses.limits.2026.current":    30000.0,
	"expenses.limits.2026.investment": 70000.0,
}
