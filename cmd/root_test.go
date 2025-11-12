package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setupTests creates a temporary directory, a mock deploy script,
// and sets the HOME env var to the temp dir.
func setupTests(t *testing.T) (tempDir string, mockDeployPath string) {
	t.Helper()

	tempDir = t.TempDir()
	t.Setenv("HOME", tempDir)

	// Create a mock Rime directory structure
	rimeDir := filepath.Join(tempDir, "Library", "Rime")
	err := os.MkdirAll(rimeDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create mock Rime directory: %v", err)
	}

	// Create a mock deploy script
	mockDeployPath = filepath.Join(tempDir, "mock_deploy.sh")
	scriptContent := "#!/bin/bash\necho 'Deployment successful!'\n"
	err = os.WriteFile(mockDeployPath, []byte(scriptContent), 0o755)
	if err != nil {
		t.Fatalf("Failed to create mock deploy script: %v", err)
	}

	return tempDir, mockDeployPath
}

// executeCommand is a helper to execute a cobra command and capture its output.
func executeCommand(t *testing.T, args ...string) (string, error) {
	t.Helper()
	b := new(bytes.Buffer)
	rootCmd.SetOut(b)
	rootCmd.SetErr(b)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return b.String(), err
}

func TestRunDeployCommand(t *testing.T) {
	tempDir, mockDeployPath := setupTests(t)
	defer os.RemoveAll(tempDir)

	// Set the deploy command to our mock script
	deployCommand = mockDeployPath

	// Redirect stdout to capture the output of the command
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runDeployCommand()

	w.Close()
	os.Stdout = oldStdout
	if err != nil {
		t.Fatalf("runDeployCommand() failed: %v", err)
	}

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "Deployment successful!") {
		t.Errorf("Expected deployment message not found in output. Got: %s", output)
	}
}

func TestListCommand(t *testing.T) {
	tempDir, _ := setupTests(t)
	defer os.RemoveAll(tempDir)
	dictPath := filepath.Join(tempDir, "Library", "Rime", "test.dict.yaml")
	content := `---
...
word1	code1	10
# comment
word2	code2	20
`
	err := os.WriteFile(dictPath, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("Failed to write temp dict file: %v", err)
	}

	userDictFile = dictPath // Override for the test

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"list"})
	rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "word1") || !strings.Contains(output, "word2") {
		t.Errorf("list output should contain the words. Got: %s", output)
	}
}

func TestAddCommand(t *testing.T) {
	tempDir, mockDeployPath := setupTests(t)
	defer os.RemoveAll(tempDir)
	userDictPath := filepath.Join(tempDir, "Library", "Rime", "user.dict.yaml")
	mainDictPath := filepath.Join(tempDir, "Library", "Rime", "main.dict.yaml")

	// Create empty user dict and a simple main dict
	os.WriteFile(userDictPath, []byte("---\n...\n"), 0o644)
	os.WriteFile(mainDictPath, []byte("测\ty\n试\tf\n"), 0o644)

	// Override flags for the test
	userDictFile = userDictPath
	mainDictFile = mainDictPath
	deployCommand = mockDeployPath

	_, err := executeCommand(t, "add", "测试", "--weight", "50", "--group", "Test")
	if err != nil {
		t.Fatalf("add command failed: %v", err)
	}

	content, _ := os.ReadFile(userDictPath)
	if !strings.Contains(string(content), "测试\tyf\t50") {
		t.Errorf("add command did not add the word correctly. File content:\n%s", content)
	}
	if !strings.Contains(string(content), "## Test") {
		t.Errorf("add command did not create the group correctly. File content:\n%s", content)
	}
}

func TestDeleteCommand(t *testing.T) {
	tempDir, mockDeployPath := setupTests(t)
	defer os.RemoveAll(tempDir)
	dictPath := filepath.Join(tempDir, "Library", "Rime", "test.dict.yaml")
	content := `---`
	content += `
...
`
	content += "word1\tcode1\t10\n"
	content += "delete_me\tdel_code\t5\n"
	content += "word2\tcode2\t20\n"
	os.WriteFile(dictPath, []byte(content), 0o644)

	userDictFile = dictPath
	deployCommand = mockDeployPath

	_, err := executeCommand(t, "delete", "delete_me")
	if err != nil {
		t.Fatalf("delete command failed: %v", err)
	}

	fileContent, _ := os.ReadFile(dictPath)
	if strings.Contains(string(fileContent), "delete_me") {
		t.Errorf("delete command did not remove the word. File content:\n%s", fileContent)
	}
}

func TestQueryCommand(t *testing.T) {
	tempDir, _ := setupTests(t)
	dictPath := filepath.Join(tempDir, "Library", "Rime", "test.dict.yaml")
	content := `---`
	content += `
...
`
	content += "find_me\tfind_code\t100\n"
	os.WriteFile(dictPath, []byte(content), 0o644)
	userDictFile = dictPath

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"query", "find_me"})
	rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "find_me") || !strings.Contains(output, "find_code") || !strings.Contains(output, "100") {
		t.Errorf("query output is incorrect. Got: %s", output)
	}
}

func TestSetWeightCommand(t *testing.T) {
	tempDir, mockDeployPath := setupTests(t)
	defer os.RemoveAll(tempDir)
	dictPath := filepath.Join(tempDir, "Library", "Rime", "test.dict.yaml")
	content := `---`
	content += `
...
`
	content += "my_word\tmy_code\t10\n"
	os.WriteFile(dictPath, []byte(content), 0o644)
	userDictFile = dictPath
	deployCommand = mockDeployPath

	_, err := executeCommand(t, "set-weight", "my_word", "999")
	if err != nil {
		t.Fatalf("set-weight command failed: %v", err)
	}

	fileContent, _ := os.ReadFile(dictPath)
	if !strings.Contains(string(fileContent), "my_word\tmy_code\t999") {
		t.Errorf("set-weight did not update the weight. File content:\n%s", fileContent)
	}
}
