package logger

import (
	"bufio"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	logFile := "test_log.log"
	defer os.Remove(logFile)
	log, err := NewLogger(logFile)
	require.NoError(t, err)
	require.NotNil(t, log)
	log.Info("this is info")
	log.Error("this is error")
	log.Done()

	file, err := os.Open(logFile)
	require.NoError(t, err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	foundInfo := false
	foundError := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "[INFO]") && strings.Contains(line, "this is info") {
			foundInfo = true
		}
		if strings.Contains(line, "[ERROR]") && strings.Contains(line, "this is error") {
			foundError = true
		}
	}

	require.True(t, foundInfo, "INFO log not found")
	require.True(t, foundError, "ERROR log not found")
}
func TestNewLogger_InvalidPath(t *testing.T) {
	_, err := NewLogger("")
	require.Error(t, err)
}
