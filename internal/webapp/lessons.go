// Package webapp derives the embedded lesson catalog and serves the web tour.
package webapp

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Lesson is one entry in the immutable, ordered lesson catalog.
type Lesson struct {
	ID       int
	Filename string
	Value    int
	Question string
	Source   string
	Hash     string
}

var lessonFileName = regexp.MustCompile(`^([0-9]{2})_[A-Za-z0-9_]+\.go$`)

var readmeRow = regexp.MustCompile("(?m)^\\|\\s*\\[`([^`]+)`\\]\\(([^)]+)\\)\\s*\\|\\s*(\\d+)\\s*\\|\\s*(.+?)\\s*\\|\\s*$")

// LoadLessons matches every embedded `[0-9][0-9]_*.go` file in files against
// exactly one README table row and returns the sorted, contiguous catalog.
// It fails hard on any gap, duplicate, or mismatch.
func LoadLessons(files fs.FS, readme []byte) ([]Lesson, error) {
	entries, err := fs.ReadDir(files, ".")
	if err != nil {
		return nil, fmt.Errorf("webapp: read lesson dir: %w", err)
	}

	type row struct {
		id       int
		value    int
		question string
	}
	rowsByFilename := make(map[string]row)
	for _, m := range readmeRow.FindAllSubmatch(readme, -1) {
		linkText, target, valueStr, question := string(m[1]), string(m[2]), string(m[3]), string(m[4])
		if linkText != target {
			continue
		}
		idMatch := lessonFileName.FindStringSubmatch(linkText)
		if idMatch == nil {
			continue
		}
		id, err := strconv.Atoi(idMatch[1])
		if err != nil {
			return nil, fmt.Errorf("webapp: README row for %q has non-numeric id: %w", linkText, err)
		}
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return nil, fmt.Errorf("webapp: README row for %q has non-integer value: %w", linkText, err)
		}
		if value < 3 || value > 5 {
			return nil, fmt.Errorf("webapp: README row for %q has value %d outside 3-5", linkText, value)
		}
		if _, dup := rowsByFilename[linkText]; dup {
			return nil, fmt.Errorf("webapp: README has duplicate row for %q", linkText)
		}
		rowsByFilename[linkText] = row{id: id, value: value, question: question}
	}

	var lessons []Lesson
	seenID := make(map[int]string)
	for _, entry := range entries {
		name := entry.Name()
		if !lessonFileName.MatchString(name) {
			continue
		}
		r, ok := rowsByFilename[name]
		if !ok {
			return nil, fmt.Errorf("webapp: lesson file %q has no matching README row", name)
		}
		if prior, dup := seenID[r.id]; dup {
			return nil, fmt.Errorf("webapp: duplicate lesson id %d (%q and %q)", r.id, prior, name)
		}
		seenID[r.id] = name
		src, err := fs.ReadFile(files, name)
		if err != nil {
			return nil, fmt.Errorf("webapp: read lesson %q: %w", name, err)
		}
		source := browserSource(string(src))
		hash := sha256.Sum256([]byte(source))
		lessons = append(lessons, Lesson{
			ID:       r.id,
			Filename: name,
			Value:    r.value,
			Question: r.question,
			Source:   source,
			Hash:     fmt.Sprintf("%x", hash[:8]),
		})
		delete(rowsByFilename, name)
	}

	for filename := range rowsByFilename {
		return nil, fmt.Errorf("webapp: README row for %q has no matching lesson file", filename)
	}

	sort.Slice(lessons, func(i, j int) bool { return lessons[i].ID < lessons[j].ID })

	for i, l := range lessons {
		want := i + 1
		if l.ID != want {
			return nil, fmt.Errorf("webapp: lesson ids must be contiguous from 1; got %d, expected %d", l.ID, want)
		}
	}

	return lessons, nil
}

func browserSource(source string) string {
	lines := strings.Split(source, "\n")
	if len(lines) == 0 || !strings.HasPrefix(lines[0], "//go:build ") {
		return source
	}
	lines = lines[1:]
	if len(lines) > 0 && lines[0] == "" {
		lines = lines[1:]
	}
	return strings.Join(lines, "\n")
}
