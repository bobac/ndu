package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bobac/ndu/pkg/ndu"
)

func analyzeRecursive(path string, config ndu.Config, currentDepth int) error {
	dirs, err := ndu.AnalyzeDir(path, config)
	if err != nil {
		return err
	}

	// Pokud je zapnutý verbose mód, vyčistíme poslední řádek
	if config.Verbose {
		fmt.Print("\r\033[K")
	}

	// Vypíšeme všechny nalezené adresáře (omezené MaxDirs)
	ndu.PrintResults(dirs, config, path)

	// Pokud jsme dosáhli maximální hloubky rekurze, končíme
	if currentDepth >= config.Recursive {
		return nil
	}

	// Omezíme počet adresářů pro rekurzivní průchod
	recursiveDirs := dirs
	if config.RecursiveDepth > 0 && len(recursiveDirs) > config.RecursiveDepth {
		recursiveDirs = recursiveDirs[:config.RecursiveDepth]
	}

	// Rekurzivně projdeme omezený počet adresářů
	for _, dir := range recursiveDirs {
		fmt.Printf("\n=> %s\n", dir.Path)
		if err := analyzeRecursive(dir.Path, config, currentDepth+1); err != nil {
			fmt.Fprintf(os.Stderr, "Chyba při analýze adresáře %s: %v\n", dir.Path, err)
			continue
		}
	}

	return nil
}

func main() {
	var config ndu.Config
	var path string
	var showHelp bool

	flag.BoolVar(&config.HumanReadable, "h", false, "Displays sizes in human readable format")
	flag.IntVar(&config.MaxDirs, "n", 0, "Number of largest directories to display")
	flag.IntVar(&config.Recursive, "r", 0, "Recursive analysis depth")
	flag.IntVar(&config.RecursiveDepth, "d", 0, "Number of directories for recursive analysis")
	flag.BoolVar(&config.Verbose, "v", false, "Shows currently processed directory")
	flag.BoolVar(&showHelp, "help", false, "Shows this help message")
	flag.Parse()

	if showHelp {
		fmt.Println("NDU - Command line utility for disk usage analysis")
		fmt.Println("(c) 2025 Robert Houser")
		fmt.Println("\nUsage:")
		fmt.Println("  ndu [switches] [directory]")
		fmt.Println("\nSwitches:")
		fmt.Println("  -h\t\tDisplays sizes in human readable format (e.g. 1.2 GB)")
		fmt.Println("  -n count\tShows only 'count' largest directories at each level")
		fmt.Println("  -r depth\tPerforms recursive analysis of directories up to 'depth' level")
		fmt.Println("  -d count\tFor 'count' largest directories at each level performs recursive analysis")
		fmt.Println("  -v\t\tShows currently processed directory")
		fmt.Println("  -help\t\tShows this help message")
		fmt.Println("\nExamples:")
		fmt.Println("  ndu -h -n 3 /")
		fmt.Println("  ndu -h -n 3 -r 2 -d 1 /")
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) > 0 {
		path = args[0]
	} else {
		path = "."
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting absolute path: %v\n", err)
		os.Exit(1)
	}

	if err := analyzeRecursive(absPath, config, 0); err != nil {
		fmt.Fprintf(os.Stderr, "Error analyzing directory: %v\n", err)
		os.Exit(1)
	}

	// If verbose mode is enabled, add a newline at the end
	if config.Verbose {
		fmt.Println()
	}
}
