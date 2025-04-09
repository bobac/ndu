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

	flag.BoolVar(&config.HumanReadable, "h", false, "Vypisuje velikosti v lidsky čitelné formě")
	flag.IntVar(&config.MaxDirs, "n", 0, "Počet největších adresářů k zobrazení")
	flag.IntVar(&config.Recursive, "r", 0, "Hloubka rekurzivního průchodu")
	flag.IntVar(&config.RecursiveDepth, "d", 0, "Počet adresářů pro rekurzivní průchod")
	flag.BoolVar(&config.Verbose, "v", false, "Zobrazuje aktuálně zpracovávaný adresář")
	flag.BoolVar(&showHelp, "help", false, "Zobrazí nápovědu")
	flag.Parse()

	if showHelp {
		fmt.Println("NDU - Command line utilita pro zjištění obsazenosti disku")
		fmt.Println("(c) 2025 Robert Houser")
		fmt.Println("\nPoužití:")
		fmt.Println("  ndu [přepínače] [adresář]")
		fmt.Println("\nPřepínače:")
		fmt.Println("  -h\t\tVypisuje velikosti v lidsky čitelné formě (např. 1.2 GB)")
		fmt.Println("  -n počet\tVypíše jen 'počet' největších adresářů na každé úrovni")
		fmt.Println("  -r hloubka\tProvede rekurzivní analýzu adresářů do úrovně 'hloubka'")
		fmt.Println("  -d počet\tPro 'počet' největších adresářů v každé úrovni provede rekurzivní analýzu")
		fmt.Println("  -v\t\tZobrazuje aktuálně zpracovávaný adresář")
		fmt.Println("  -help\t\tZobrazí tuto nápovědu")
		fmt.Println("\nPříklady:")
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
		fmt.Fprintf(os.Stderr, "Chyba při získávání absolutní cesty: %v\n", err)
		os.Exit(1)
	}

	if err := analyzeRecursive(absPath, config, 0); err != nil {
		fmt.Fprintf(os.Stderr, "Chyba při analýze adresáře: %v\n", err)
		os.Exit(1)
	}

	// Pokud je zapnutý verbose mód, přidáme na konec nový řádek
	if config.Verbose {
		fmt.Println()
	}
}
