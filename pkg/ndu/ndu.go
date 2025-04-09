package ndu

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/term"
)

type DirSize struct {
	Path string
	Size int64
}

type Config struct {
	HumanReadable  bool
	MaxDepth       int
	MaxDirs        int
	Recursive      int
	RecursiveDepth int
	Verbose        bool
}

func FormatSize(size int64, humanReadable bool) string {
	if !humanReadable {
		return fmt.Sprintf("%d", size)
	}

	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func getTerminalWidth() int {
	// Výchozí šířka terminálu
	width := 80

	// Zkusíme získat šířku terminálu
	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		return width
	}

	return width
}

func shortenPath(path string) string {
	width := getTerminalWidth()
	// Necháme místo pro "Zpracovávám: " a nějaké rezervy
	maxLen := width - 20

	if len(path) <= maxLen {
		return path
	}

	// Zkrátíme cestu uprostřed
	half := maxLen / 2
	return path[:half-3] + "..." + path[len(path)-(half-3):]
}

func GetDirSize(path string, config Config) (int64, error) {
	var size int64
	var currentFile string
	var lastProgress string
	var lastDir string

	_ = filepath.WalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			currentFile = filePath
			if config.Verbose {
				currentDir := filepath.Dir(filePath)
				if currentDir != lastDir || len(currentFile) < len(lastProgress) {
					progress := fmt.Sprintf("\r\033[KZpracovávám: %s", shortenPath(currentFile))
					fmt.Print(progress)
					lastProgress = progress
					lastDir = currentDir
				}
			}
			info, err := d.Info()
			if err != nil {
				return nil
			}
			size += info.Size()
		}
		return nil
	})
	return size, nil
}

func AnalyzeDir(path string, config Config) ([]DirSize, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return []DirSize{}, nil
	}

	var dirs []DirSize
	for _, entry := range entries {
		if entry.IsDir() {
			dirPath := filepath.Join(path, entry.Name())
			if config.Verbose {
				fmt.Printf("\r\033[KZpracovávám adresář: %s", shortenPath(dirPath))
			}
			size, err := GetDirSize(dirPath, config)
			if err != nil {
				continue
			}
			dirs = append(dirs, DirSize{Path: dirPath, Size: size})
		}
	}

	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Size > dirs[j].Size
	})

	if config.MaxDirs > 0 && len(dirs) > config.MaxDirs {
		dirs = dirs[:config.MaxDirs]
	}

	return dirs, nil
}

func PrintResults(dirs []DirSize, config Config, prefix string) {
	if config.Verbose {
		fmt.Print("\r\033[K")
	}

	if config.HumanReadable {
		// Najdeme maximální délku cesty a velikosti
		maxPathLen := 0
		maxSizeLen := 0
		for _, dir := range dirs {
			relPath := strings.TrimPrefix(dir.Path, prefix)
			if relPath == "" {
				relPath = "."
			}
			if len(relPath) > maxPathLen {
				maxPathLen = len(relPath)
			}
			sizeStr := FormatSize(dir.Size, true)
			if len(sizeStr) > maxSizeLen {
				maxSizeLen = len(sizeStr)
			}
		}

		// Vypíšeme výsledky s zarovnáním
		for _, dir := range dirs {
			relPath := strings.TrimPrefix(dir.Path, prefix)
			if relPath == "" {
				relPath = "."
			}
			sizeStr := FormatSize(dir.Size, true)
			fmt.Printf("%-*s  %*s\n", maxPathLen, relPath, maxSizeLen, sizeStr)
		}
	} else {
		// Vypíšeme výsledky bez zarovnání
		for _, dir := range dirs {
			relPath := strings.TrimPrefix(dir.Path, prefix)
			if relPath == "" {
				relPath = "."
			}
			fmt.Printf("%s\t%d\n", relPath, dir.Size)
		}
	}
}
