package server

import (
	"os"
	"path/filepath"
	"strconv"
)

func getPidFilePath() string {
	return filepath.Join(os.TempDir(), "espigol_server.pid")
}

func writePidFile(pid int) error {
	path := getPidFilePath()
	return os.WriteFile(path, []byte(strconv.Itoa(pid)), 0644)
}

func readPidFile() (int, error) {
	path := getPidFilePath()
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(data))
}

func removePidFile() error {
	path := getPidFilePath()
	return os.Remove(path)
}
