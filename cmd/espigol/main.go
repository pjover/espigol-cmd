package main

import "github.com/pjover/espigol/internal"

func main() {
	internal.InjectDependencies().Execute()
}
