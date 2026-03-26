package dynamicchecks_test

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gpma/framework"
	_ "gpma/framework/plugin/dynamic_checks"
)

// Run with -update to regenerate golden files:
//
//	go test ./plugin/dynamic_checks/... -update
var update = flag.Bool("update", false, "update golden files instead of comparing")

const (
	samplesDir = "../../../samples/projects"
	goldenDir  = "testdata/golden"
)

func TestDynamicChecksGolden(t *testing.T) {
	projectEntries, err := os.ReadDir(samplesDir)
	if err != nil {
		t.Fatalf("cannot read samples directory %q: %v", samplesDir, err)
	}

	for _, projectEntry := range projectEntries {
		if !projectEntry.IsDir() {
			continue
		}

		projectName := projectEntry.Name()
		projectPath := filepath.Join(samplesDir, projectName)

		goFiles, err := filepath.Glob(filepath.Join(projectPath, "*.go"))
		if err != nil {
			t.Fatalf("glob failed for project %q: %v", projectName, err)
		}

		for _, srcPath := range goFiles {
			fileName := filepath.Base(srcPath)

			goldenPath := filepath.Join(goldenDir, projectName, fileName+".golden")

			testName := projectName + "/" + fileName
			t.Run(testName, func(t *testing.T) {
				got, err := runDynamicChecksPlugin(srcPath)
				if err != nil {
					t.Fatalf("plugin processing failed for %q: %v", srcPath, err)
				}

				if *update {
					writeGolden(t, goldenPath, got)
					return
				}

				want := readGolden(t, goldenPath)
				if normalise(got) != normalise(want) {
					t.Errorf("output mismatch for %q\n\n--- want ---\n%s\n--- got ---\n%s",
						srcPath, want, got)
				}
			})
		}
	}
}

// runDynamicChecksPlugin runs the full framework pipeline (parse → bind → plugin → codegen)
// on the given source file using the dynamic_checks plugin and returns the generated code.
func runDynamicChecksPlugin(srcPath string) (string, error) {
	fw := framework.New()
	fw, err := fw.AddPluginByName("dynamic_checks")
	if err != nil {
		return "", err
	}

	result, err := fw.ProcessFile(srcPath)
	if err != nil {
		return "", err
	}

	if result.Modified {
		return result.GeneratedCode, nil
	}

	// No modifications — return the original source unchanged.
	raw, err := os.ReadFile(srcPath)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func writeGolden(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("cannot create golden directory for %q: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("cannot write golden file %q: %v", path, err)
	}
	t.Logf("updated golden file: %s", path)
}

func readGolden(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("golden file %q not found — run with -update to generate it: %v", path, err)
	}
	return string(data)
}

func normalise(s string) string {
	lines := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	for i, l := range lines {
		lines[i] = strings.TrimRight(l, " \t")
	}
	return strings.TrimRight(strings.Join(lines, "\n"), "\n")
}
