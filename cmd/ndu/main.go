package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/bobac/ndu/pkg/ndu"
)

func analyzeRecursive(path string, config ndu.Config, currentDepth int) (ndu.JSONDir, error) {
	dirs, err := ndu.AnalyzeDir(path, config)
	if err != nil {
		return ndu.JSONDir{}, err
	}

	// Pokud je zapnutý verbose mód, vyčistíme poslední řádek
	if config.Verbose {
		fmt.Print("\r\033[K")
	}

	// Vypíšeme všechny nalezené adresáře (omezené MaxDirs)
	ndu.PrintResults(dirs, config, path)

	// Vytvoříme JSON strukturu pro aktuální úroveň
	jsonDir, err := ndu.ExportToJSON(dirs, config, path)
	if err != nil {
		return ndu.JSONDir{}, err
	}

	// Pokud jsme dosáhli maximální hloubky rekurze, končíme
	if currentDepth >= config.Recursive {
		return jsonDir, nil
	}

	// Omezíme počet adresářů pro rekurzivní průchod
	recursiveDirs := dirs
	if config.RecursiveDepth > 0 && len(recursiveDirs) > config.RecursiveDepth {
		recursiveDirs = recursiveDirs[:config.RecursiveDepth]
	}

	// Rekurzivně projdeme omezený počet adresářů
	for i, dir := range recursiveDirs {
		fmt.Printf("\n=> %s\n", dir.Path)
		childJSON, err := analyzeRecursive(dir.Path, config, currentDepth+1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Chyba při analýze adresáře %s: %v\n", dir.Path, err)
			continue
		}
		jsonDir.Children[i].Children = childJSON.Children
	}

	return jsonDir, nil
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
	flag.StringVar(&config.JSONOutput, "j", "", "Export results to JSON file")
	flag.StringVar(&config.HTMLOutput, "html", "", "Export results to HTML file with pie chart visualization")
	flag.BoolVar(&showHelp, "help", false, "Shows this help message")
	autoMode := flag.Bool("auto", false, "Auto mode with reasonable defaults and opens HTML in default browser")
	flag.Parse()

	if *autoMode {
		// Nastavíme výchozí parametry pro auto mód
		config.HumanReadable = true
		config.MaxDirs = 10
		config.Recursive = 4
		config.RecursiveDepth = 5
		config.Verbose = true
		config.HTMLOutput = "auto.html"
	}

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
		fmt.Println("  -j file.json\tExport results to JSON file")
		fmt.Println("  -html file.html\tExport results to HTML file with pie chart visualization")
		fmt.Println("  -help\t\tShows this help message")
		fmt.Println("  -auto\t\tAuto mode with reasonable defaults and opens HTML in default browser")
		fmt.Println("\nExamples:")
		fmt.Println("  ndu -h -n 3 /")
		fmt.Println("  ndu -h -n 3 -r 2 -d 1 /")
		fmt.Println("  ndu -j results.json /")
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

	jsonResult, err := analyzeRecursive(absPath, config, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error analyzing directory: %v\n", err)
		os.Exit(1)
	}

	// If JSON output is requested, save the results
	if config.JSONOutput != "" {
		jsonData, err := json.MarshalIndent(jsonResult, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
			os.Exit(1)
		}

		if err := os.WriteFile(config.JSONOutput, jsonData, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing JSON file: %v\n", err)
			os.Exit(1)
		}
	}

	// If HTML output is requested, save the results
	if config.HTMLOutput != "" {
		htmlData := ndu.ExportToHTML(jsonResult)
		if err := os.WriteFile(config.HTMLOutput, []byte(htmlData), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing HTML file: %v\n", err)
			os.Exit(1)
		}

		// Pokud jsme v auto módu, otevřeme HTML v prohlížeči
		if *autoMode {
			var cmd *exec.Cmd
			switch runtime.GOOS {
			case "windows":
				cmd = exec.Command("cmd", "/c", "start", config.HTMLOutput)
			case "darwin":
				cmd = exec.Command("open", config.HTMLOutput)
			default: // linux a ostatní
				cmd = exec.Command("xdg-open", config.HTMLOutput)
			}
			if err := cmd.Start(); err != nil {
				fmt.Fprintf(os.Stderr, "Error opening browser: %v\n", err)
			}
		}
	}

	// If verbose mode is enabled, add a newline at the end
	if config.Verbose {
		fmt.Println()
	}
}
