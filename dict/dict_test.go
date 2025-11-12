package dict

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// Helper function to create a temporary dictionary file for testing.
func createTempDictFile(t *testing.T, dir, content string) string {
	t.Helper()
	path := filepath.Join(dir, "test.dict.yaml")
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp dict file: %v", err)
	}
	return path
}

func TestDictionary_Load(t *testing.T) {
	content := `# Rime dictionary
#
---
name: wubi86_jidian_user
version: "1.0"
...
一丁	ag
丁一	sa
## Custom Phrases
丁丁	ss
`
	tempDir := t.TempDir()
	dictPath := createTempDictFile(t, tempDir, content)

	d := NewDictionary(dictPath)
	err := d.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if len(d.Header) != 6 {
		t.Errorf("Expected 6 header lines, got %d", len(d.Header))
	}

	if len(d.Entries) != 4 {
		t.Errorf("Expected 4 entries, got %d", len(d.Entries))
	}

	expectedEntries := []Entry{
		{Word: "一丁", Code: "ag", Weight: 0, RawLine: "一丁\tag"},
		{Word: "丁一", Code: "sa", Weight: 0, RawLine: "丁一\tsa"},
		{IsGroup: true, Group: "Custom Phrases", RawLine: "## Custom Phrases"},
		{Word: "丁丁", Code: "ss", Weight: 0, RawLine: "丁丁\tss"},
	}

	for i, entry := range d.Entries {
		// Don't compare RawLine for simplicity, as it's tested implicitly by the other fields.
		entry.RawLine = ""
		expectedEntries[i].RawLine = ""
		if !reflect.DeepEqual(entry, expectedEntries[i]) {
			t.Errorf("Entry %d mismatch:\ngot:  %+v\nwant: %+v", i, entry, expectedEntries[i])
		}
	}
}

func TestDictionary_Save(t *testing.T) {
	tempDir := t.TempDir()
	dictPath := filepath.Join(tempDir, "save_test.dict.yaml")

	d := &Dictionary{
		path:   dictPath,
		Header: []string{"---", "..."},
		Entries: []Entry{
			{Word: "测试", Code: "iyf", Weight: 1},
			{IsGroup: true, Group: "My Group"},
			{Word: "词语", Code: "yiy", Weight: 0},
			{IsComment: true, Comment: "# A comment"},
		},
	}

	err := d.Save()
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	content, err := os.ReadFile(dictPath)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	expectedContent := `---`
expectedContent += `
...`
expectedContent += `
测试	iyf	1`
expectedContent += `
## My Group`
expectedContent += `
词语	yiy	0`
expectedContent += `
# A comment`
expectedContent += "\n"
	if string(content) != expectedContent {
		t.Errorf("Saved content mismatch:\ngot:\n%s\nwant:\n%s", string(content), expectedContent)
	}
}

func TestDictionary_AddOrUpdate(t *testing.T) {
	// Test updating an existing word
	d := &Dictionary{
		Entries: []Entry{
			{Word: "word1", Code: "c1", Weight: 1},
		},
	}
	d.AddOrUpdate("word1", "new_code", 100, "group1")
	if d.Entries[0].Code != "new_code" || d.Entries[0].Weight != 100 {
		t.Errorf("Failed to update existing word. Got: %+v", d.Entries[0])
	}

	// Test adding a new word to an existing group
	d = &Dictionary{
		Entries: []Entry{
			{IsGroup: true, Group: "group1"},
			{Word: "word1", Code: "c1", Weight: 1},
		},
	}
	d.AddOrUpdate("word2", "c2", 2, "group1")
	if len(d.Entries) != 3 || d.Entries[1].Word != "word2" {
		t.Errorf("Failed to add new word to existing group. Entries: %+v", d.Entries)
	}

	// Test adding a new word and creating a new group
	d = &Dictionary{
		Entries: []Entry{
			{Word: "word1", Code: "c1", Weight: 1},
		},
	}
	d.AddOrUpdate("word2", "c2", 2, "new_group")
	if len(d.Entries) != 3 || !d.Entries[1].IsGroup || d.Entries[1].Group != "new_group" || d.Entries[2].Word != "word2" {
		t.Errorf("Failed to add new word and create new group. Entries: %+v", d.Entries)
	}
}

func TestWubiEncoder_GenerateCode(t *testing.T) {
	mainDictContent := `中	k
国	l
人	w
民	n
`
	tempDir := t.TempDir()
	mainDictPath := createTempDictFile(t, tempDir, mainDictContent)

	encoder, err := NewWubiEncoder(mainDictPath)
	if err != nil {
		t.Fatalf("NewWubiEncoder failed: %v", err)
	}

	testCases := []struct {
		word     string
		expected string
		hasError bool
	}{
		{"中", "k", false},
		{"中国", "kl", false},
		{"中国人", "klw", false},
		{"中国人民", "klwn", false},
		{"测试", "", true}, // "测" is not in the dictionary
		{"", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.word, func(t *testing.T) {
			code, err := encoder.GenerateCode(tc.word)
			if (err != nil) != tc.hasError {
				t.Errorf("Expected error: %v, got: %v", tc.hasError, err)
			}
			if code != tc.expected {
				t.Errorf("Expected code: %s, got: %s", tc.expected, code)
			}
		})
	}
}
