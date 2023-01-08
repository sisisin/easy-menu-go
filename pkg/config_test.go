package pkg

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestLoadConfigNormally(t *testing.T) {
	cleanup := mockGetwd(t, "case_exists/d1")
	defer cleanup()
	cleanup = mockExitWithPathError(t, nil)
	defer cleanup()

	configPath, document := LoadConfig("")

	assert.Equal(t, path.Join(getFixturesPath(t), "case_exists/d1/easy-menu.yaml"), configPath)
	assert.NotEqual(t, yaml.Node{}, *document)
}

func TestLoadConfigNotFound(t *testing.T) {
	cleanup := mockGetwd(t, "case_not_exists/d1")
	defer cleanup()

	defer func() {
		err := recover()
		if err != "from test code" {
			t.Fatal(err)
		}
	}()

	exitCode := 0
	errorMsg := ""
	cleanup = mockExitWithPathError(t, func(code int, msg string) {
		exitCode = code
		errorMsg = msg
		panic("from test code")
	})
	defer cleanup()

	LoadConfig("")
	assert.Equal(t, exitCode, 1)
	assert.Equal(t, errorMsg, "")
}

func mockExitWithPathError(t *testing.T, mocked func(code int, msg string)) func() {
	old := exit
	cleanup := func() {
		exit = old
	}

	if mocked == nil {
		exit = func(code int, msg string) {}
	} else {
		exit = mocked
	}

	return cleanup
}

func mockGetwd(t *testing.T, dir string) func() {
	old := getwd
	cleanup := func() {
		getwd = old
	}

	getwd = func() (string, error) {
		return path.Join(getFixturesPath(t), dir), nil
	}
	return cleanup
}

func getFixturesPath(t *testing.T) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	return path.Join(wd, "fixtures/config_cases")
}
