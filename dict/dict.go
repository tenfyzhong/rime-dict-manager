package dict

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Entry represents a single line in the dictionary file.
type Entry struct {
	Word      string
	Code      string
	Weight    int
	Comment   string // For standalone comment lines
	Group     string // For group lines like '## GroupName'
	IsComment bool   // True if the line is a comment
	IsGroup   bool   // True if the line is a group header
	RawLine   string // The original, unmodified line
}

// Dictionary holds the entire content of a dictionary file.
type Dictionary struct {
	Header  []string // YAML header part
	Entries []Entry
	path    string
}

// NewDictionary creates a new Dictionary instance.
func NewDictionary(path string) *Dictionary {
	return &Dictionary{path: path}
}

// Load reads and parses the dictionary file.
func (d *Dictionary) Load() error {
	file, err := os.Open(d.path)
	if err != nil {
		return fmt.Errorf("failed to open dictionary file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inHeader := true
	for scanner.Scan() {
		line := scanner.Text()

		if inHeader {
			if strings.HasPrefix(line, "...") {
				inHeader = false
			}
			d.Header = append(d.Header, line)
			continue
		}

		var entry Entry
		entry.RawLine = line

		if strings.HasPrefix(line, "##") {
			entry.IsGroup = true
			entry.Group = strings.TrimSpace(strings.TrimPrefix(line, "##"))
		} else if strings.HasPrefix(line, "#") {
			entry.IsComment = true
			entry.Comment = line
		} else if strings.TrimSpace(line) != "" {
			parts := strings.Split(line, "\t")
			if len(parts) >= 2 {
				entry.Word = parts[0]
				entry.Code = parts[1]
				if len(parts) > 2 {
					weight, err := strconv.Atoi(parts[2])
					if err == nil {
						entry.Weight = weight
					}
				}
			}
		}
		d.Entries = append(d.Entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading dictionary file: %w", err)
	}

	return nil
}

// Save writes the dictionary content back to the file.
func (d *Dictionary) Save() error {
	file, err := os.Create(d.path)
	if err != nil {
		return fmt.Errorf("failed to create dictionary file for writing: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, line := range d.Header {
		_, _ = writer.WriteString(line + "\n")
	}

	for _, entry := range d.Entries {
		if entry.IsGroup {
			_, _ = writer.WriteString(fmt.Sprintf("## %s\n", entry.Group))
		} else if entry.IsComment {
			_, _ = writer.WriteString(entry.Comment + "\n")
		} else if entry.Word != "" {
			_, _ = writer.WriteString(fmt.Sprintf("%s	%s	%d\n", entry.Word, entry.Code, entry.Weight))
		} else {
			_, _ = writer.WriteString(entry.RawLine + "\n")
		}
	}

	return writer.Flush()
}

// AddOrUpdate finds a word and updates it, or adds it if it doesn't exist.
func (d *Dictionary) AddOrUpdate(word, code string, weight int, group string) {
	// First, try to update existing entry
	for i := range d.Entries {
		if d.Entries[i].Word == word {
			d.Entries[i].Code = code
			d.Entries[i].Weight = weight
			// Note: Moving an entry to a different group is complex.
			// For now, we just update it in place.
			return
		}
	}

	// If not found, add a new entry
	newEntry := Entry{
		Word:   word,
		Code:   code,
		Weight: weight,
	}

	// Find the target group and insert the new entry
	groupFound := false
	for i, entry := range d.Entries {
		if entry.IsGroup && entry.Group == group {
			// Insert after the group header
			d.Entries = append(d.Entries[:i+1], append([]Entry{newEntry}, d.Entries[i+1:]...)...)
			groupFound = true
			break
		}
	}

	// If group is not found, create it at the end of the file
	if !groupFound {
		d.Entries = append(d.Entries, Entry{IsGroup: true, Group: group})
		d.Entries = append(d.Entries, newEntry)
	}
}

// WubiEncoder can generate Wubi codes for Chinese words.
type WubiEncoder struct {
	charMap map[rune]string
}

// NewWubiEncoder creates an encoder by loading a main dictionary file.
func NewWubiEncoder(mainDictPath string) (*WubiEncoder, error) {
	charMap := make(map[rune]string)

	file, err := os.Open(mainDictPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open main dictionary '%s': %w", mainDictPath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) >= 2 {
			char := []rune(parts[0])
			if len(char) == 1 {
				// We only care about single characters for building words
				charMap[char[0]] = parts[1]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading main dictionary: %w", err)
	}

	return &WubiEncoder{charMap: charMap}, nil
}

// GenerateCode generates a Wubi code for a given word.
// Rules:
// 1-char word: full code (up to 4 letters)
// 2-char word: 1st letter of 1st char + 1st letter of 2nd char
// 3-char word: 1st letter of 1st, 2nd, 3rd chars
// 4+ char word: 1st letter of 1st, 2nd, 3rd chars + 1st letter of last char
func (e *WubiEncoder) GenerateCode(word string) (string, error) {
	runes := []rune(word)
	if len(runes) == 0 {
		return "", nil
	}

	if len(runes) == 1 {
		code, ok := e.charMap[runes[0]]
		if !ok {
			return "", fmt.Errorf("character '%c' not found in main dictionary", runes[0])
		}
		return code, nil
	}

	var generatedCode strings.Builder
	var codesToTake []rune

	if len(runes) <= 3 {
		codesToTake = runes
	} else { // 4 or more characters
		codesToTake = []rune{runes[0], runes[1], runes[2], runes[len(runes)-1]}
	}

	for _, r := range codesToTake {
		code, ok := e.charMap[r]
		if !ok {
			return "", fmt.Errorf("character '%c' not found in main dictionary", r)
		}
		if len(code) > 0 {
			generatedCode.WriteByte(code[0])
		}
	}

	return generatedCode.String(), nil
}
