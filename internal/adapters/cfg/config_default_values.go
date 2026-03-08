package cfg

var defaultValues = map[string]interface{}{
	"business.name":                   "Cooperativa d'Estellencs",
	"db.server":                       "mongodb://localhost:27017",
	"db.name":                         "espigol",
	"expenses.limits.2026.current":    30000.0,
	"expenses.limits.2026.investment": 70000.0,
	"files.logo":                      "configs/logo.png",
	"output.directory":                "~/espigol/reports",
	"server.port":                     "8080",
}
