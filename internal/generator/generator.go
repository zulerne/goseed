package generator

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/zulerne/goseed/internal/config"
)

const (
	dirPerm  = 0o750
	filePerm = 0o600
)

// Generate creates a project from the given config using templates from fsys.
func Generate(cfg *config.ProjectConfig, fsys fs.FS, outputDir string) error {
	targetDir := filepath.Join(outputDir, cfg.ProjectName)
	if err := os.MkdirAll(targetDir, dirPerm); err != nil {
		return fmt.Errorf("create project dir: %w", err)
	}

	var count int

	for _, fm := range Manifest {
		if fm.Condition != nil && !fm.Condition(cfg) {
			continue
		}

		target, err := resolveTarget(fm.Target, cfg)
		if err != nil {
			return fmt.Errorf("resolve target %q: %w", fm.Target, err)
		}

		destPath := filepath.Join(targetDir, target)

		err = os.MkdirAll(filepath.Dir(destPath), dirPerm)
		if err != nil {
			return fmt.Errorf("create dir for %q: %w", target, err)
		}

		srcPath := "templates/" + fm.Source
		data, err := fs.ReadFile(fsys, srcPath)
		if err != nil {
			return fmt.Errorf("read template %q: %w", srcPath, err)
		}

		if fm.IsTemplate {
			data, err = renderTemplate(fm.Source, string(data), cfg)
			if err != nil {
				return fmt.Errorf("render %q: %w", fm.Source, err)
			}
		}

		err = os.WriteFile(destPath, data, filePerm)
		if err != nil {
			return fmt.Errorf("write %q: %w", target, err)
		}

		count++
	}

	fmt.Printf("  Created %s/\n", cfg.ProjectName)
	fmt.Printf("  %d files generated\n", count)

	runPostGen(targetDir)

	return nil
}

func resolveTarget(target string, cfg *config.ProjectConfig) (string, error) {
	tmpl, err := template.New("target").Parse(target)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cfg); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func renderTemplate(name, content string, cfg *config.ProjectConfig) ([]byte, error) {
	tmpl, err := template.New(name).Parse(content)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cfg); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func runPostGen(dir string) {
	if _, err := exec.LookPath("go"); err == nil {
		cmd := exec.Command("go", "mod", "tidy")
		cmd.Dir = dir
		if err := cmd.Run(); err == nil {
			fmt.Println("  go mod tidy done")
		}
	}

	if _, err := exec.LookPath("git"); err == nil {
		cmd := exec.Command("git", "init")
		cmd.Dir = dir
		cmd.Stdout = nil
		cmd.Stderr = nil
		if err := cmd.Run(); err == nil {
			add := exec.Command("git", "add", ".")
			add.Dir = dir
			if err := add.Run(); err == nil {
				fmt.Println("  git init done")
			}
		}
	}

	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", filepath.Base(dir))
	fmt.Printf("  go run ./cmd/%s\n", filepath.Base(dir))
}
