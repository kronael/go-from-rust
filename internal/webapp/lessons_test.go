package webapp

import (
	"strings"
	"testing"
	"testing/fstest"
)

const readmeHeader = "| Lesson | Value | Main point |\n|---|:---:|---|\n"

func row(filename string, value int, summary string) string {
	return "| [`" + filename + "`](" + filename + ") | " + itoa(value) + " | " + summary + " |\n"
}

func itoa(n int) string {
	return string(rune('0' + n))
}

func mapFS(filenames ...string) fstest.MapFS {
	files := fstest.MapFS{}
	for _, name := range filenames {
		files[name] = &fstest.MapFile{Data: []byte("package main\n")}
	}
	return files
}

func TestLoadLessonsSuccess(t *testing.T) {
	files := mapFS("01_arrays_slices.go", "02_copy.go", "03_range.go")
	files["01_arrays_slices.go"] = &fstest.MapFile{Data: []byte("//go:build ignore\n\npackage main\n")}
	readme := readmeHeader +
		row("01_arrays_slices.go", 5, "Assigned slices share storage.") +
		row("02_copy.go", 5, "Clone creates; copy fills.") +
		row("03_range.go", 5, "range yields operand-specific values.")

	lessons, err := LoadLessons(files, []byte(readme))
	if err != nil {
		t.Fatalf("LoadLessons: %v", err)
	}
	if len(lessons) != 3 {
		t.Fatalf("got %d lessons, want 3", len(lessons))
	}
	for i, l := range lessons {
		if l.ID != i+1 {
			t.Errorf("lessons[%d].ID = %d, want %d", i, l.ID, i+1)
		}
	}
	if lessons[0].Value != 5 || lessons[0].Summary != "Assigned slices share storage." {
		t.Errorf("lessons[0] = %+v", lessons[0])
	}
	if !strings.HasPrefix(lessons[0].Source, "package main") {
		t.Errorf("lessons[0].Source = %q", lessons[0].Source)
	}
	if lessons[0].Hash == "" {
		t.Error("lessons[0].Hash is empty")
	}
	if strings.Contains(lessons[0].Source, "go:build") {
		t.Errorf("lessons[0].Source contains repository build tag: %q", lessons[0].Source)
	}
}

func TestLoadLessonsRejectsValueOutsideEditorialRange(t *testing.T) {
	files := mapFS("01_arrays_slices.go")
	readme := readmeHeader + row("01_arrays_slices.go", 2, "Q1?")

	_, err := LoadLessons(files, []byte(readme))
	if err == nil {
		t.Fatal("expected error for value outside 3-5, got nil")
	}
}

func TestBrowserSourceOnlyRemovesLeadingBuildConstraint(t *testing.T) {
	source := "//go:build ignore\n\npackage main\n\n//go:build stays\n"
	want := "package main\n\n//go:build stays\n"
	if got := browserSource(source); got != want {
		t.Fatalf("browserSource() = %q, want %q", got, want)
	}
}

func TestLoadLessonsGap(t *testing.T) {
	files := mapFS("01_arrays_slices.go", "03_range.go")
	readme := readmeHeader +
		row("01_arrays_slices.go", 5, "Q1?") +
		row("03_range.go", 5, "Q3?")

	_, err := LoadLessons(files, []byte(readme))
	if err == nil {
		t.Fatal("expected error for id gap, got nil")
	}
}

func TestLoadLessonsDuplicateID(t *testing.T) {
	files := fstest.MapFS{
		"01_arrays_slices.go": &fstest.MapFile{Data: []byte("package main\n")},
		"01_other.go":         &fstest.MapFile{Data: []byte("package main\n")},
	}
	readme := readmeHeader +
		row("01_arrays_slices.go", 5, "Q1?") +
		row("01_other.go", 5, "Q1 again?")

	_, err := LoadLessons(files, []byte(readme))
	if err == nil {
		t.Fatal("expected error for duplicate id, got nil")
	}
}

func TestLoadLessonsMismatchedLink(t *testing.T) {
	files := mapFS("01_arrays_slices.go")
	readme := readmeHeader +
		"| [`01_arrays_slices.go`](01_other_target.go) | 5 | Q1? |\n"

	_, err := LoadLessons(files, []byte(readme))
	if err == nil {
		t.Fatal("expected error for mismatched link text/target, got nil")
	}
}

func TestLoadLessonsUnmatchedFile(t *testing.T) {
	files := mapFS("01_arrays_slices.go", "02_copy.go")
	readme := readmeHeader +
		row("01_arrays_slices.go", 5, "Q1?")

	_, err := LoadLessons(files, []byte(readme))
	if err == nil {
		t.Fatal("expected error for unmatched lesson file, got nil")
	}
}

func TestLoadLessonsUnmatchedRow(t *testing.T) {
	files := mapFS("01_arrays_slices.go")
	readme := readmeHeader +
		row("01_arrays_slices.go", 5, "Q1?") +
		row("02_copy.go", 5, "Q2?")

	_, err := LoadLessons(files, []byte(readme))
	if err == nil {
		t.Fatal("expected error for unmatched README row, got nil")
	}
}
