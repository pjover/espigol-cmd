// Package main is the entry point of the espigol application.
//
// @title           Espigol API
// @version         1.0
// @description     REST API for managing partners and expense forecasts.
// @contact.name    Espigol
// @license.name    MIT
// @host            localhost:8080
// @BasePath        /
package main

import (
	_ "github.com/pjover/espigol/docs"
	"github.com/pjover/espigol/internal"
)

func main() {
	internal.InjectDependencies().Execute()
}
