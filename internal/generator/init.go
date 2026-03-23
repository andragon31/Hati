package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

func InitProject(root string) error {
	hatiDir := filepath.Join(root, ".hati")
	if err := os.MkdirAll(hatiDir, 0755); err != nil {
		return err
	}
	fmt.Printf("Hati initialized in %s\n", hatiDir)
	return nil
}
