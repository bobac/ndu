package ndu

import (
	"encoding/json"
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
	JSONOutput     string
	HTMLOutput     string
}

type JSONDir struct {
	Path     string    `json:"path"`
	Size     int64     `json:"size"`
	Children []JSONDir `json:"children,omitempty"`
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
	// Default terminal width
	width := 80

	// Try to get terminal width
	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		return width
	}

	return width
}

func shortenPath(path string) string {
	width := getTerminalWidth()
	// Leave space for "Processing: " and some margin
	maxLen := width - 20

	if len(path) <= maxLen {
		return path
	}

	// Shorten path in the middle
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
					progress := fmt.Sprintf("\r\033[KProcessing: %s", shortenPath(currentFile))
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
				fmt.Printf("\r\033[KProcessing directory: %s", shortenPath(dirPath))
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
		// Find maximum path and size lengths
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

		// Print results with alignment
		for _, dir := range dirs {
			relPath := strings.TrimPrefix(dir.Path, prefix)
			if relPath == "" {
				relPath = "."
			}
			sizeStr := FormatSize(dir.Size, true)
			fmt.Printf("%-*s  %*s\n", maxPathLen, relPath, maxSizeLen, sizeStr)
		}
	} else {
		// Print results without alignment
		for _, dir := range dirs {
			relPath := strings.TrimPrefix(dir.Path, prefix)
			if relPath == "" {
				relPath = "."
			}
			fmt.Printf("%s\t%d\n", relPath, dir.Size)
		}
	}
}

func ExportToJSON(dirs []DirSize, config Config, prefix string) (JSONDir, error) {
	root := JSONDir{
		Path: prefix,
		Size: 0,
	}

	for _, dir := range dirs {
		relPath := strings.TrimPrefix(dir.Path, prefix)
		if relPath == "" {
			relPath = "."
		}
		root.Size += dir.Size
		root.Children = append(root.Children, JSONDir{
			Path: dir.Path,
			Size: dir.Size,
		})
	}

	return root, nil
}

func ExportToHTML(jsonDir JSONDir) string {
	// Převedeme data na JSON
	jsonData, err := json.Marshal(jsonDir)
	if err != nil {
		return fmt.Sprintf("Error generating JSON: %v", err)
	}

	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Disk Usage Analysis</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
            background-color: #f5f5f5;
        }
        .header {
            display: flex;
            align-items: center;
            margin-bottom: 20px;
            background-color: white;
            padding: 15px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .up-btn {
            padding: 8px 15px;
            background-color: #2196F3;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            margin-right: 10px;
            font-weight: bold;
        }
        .up-btn:hover {
            background-color: #1976D2;
        }
        .up-btn:disabled {
            background-color: #BDBDBD;
            cursor: not-allowed;
        }
        .path {
            flex-grow: 1;
            font-family: monospace;
            padding: 8px;
            background-color: #f0f0f0;
            border-radius: 4px;
            margin-right: 10px;
        }
        .size {
            font-weight: bold;
            margin-right: 10px;
        }
        .copy-btn {
            padding: 8px 15px;
            background-color: #4CAF50;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
        }
        .copy-btn:hover {
            background-color: #45a049;
        }
        .chart-container {
            background-color: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            max-width: 800px;
            max-height: 600px;
            margin: 0 auto;
        }
        canvas {
            max-height: 500px !important;
        }
    </style>
</head>
<body>
    <div class="header">
        <button class="up-btn" onclick="goUp()" id="upButton">↑</button>
        <div class="path" id="currentPath"></div>
        <div class="size" id="currentSize"></div>
        <button class="copy-btn" onclick="copyPath()">Copy Path</button>
    </div>
    <div class="chart-container">
        <canvas id="pieChart"></canvas>
    </div>

    <script>
        let currentData = ` + string(jsonData) + `;
        let pathHistory = [];
        let chart;

        // Funkce pro generování barev pro adresáře s dětmi
        function generateColor(index, total) {
            const colors = [
                '#4CAF50', // zelená
                '#2196F3', // modrá
                '#9C27B0', // fialová
                '#FF9800', // oranžová
                '#E91E63', // růžová
                '#009688', // tyrkysová
                '#673AB7', // tmavě fialová
                '#3F51B5', // indigo
                '#00BCD4', // světle modrá
                '#8BC34A'  // světle zelená
            ];
            return colors[index % colors.length];
        }

        // Funkce pro generování odstínů šedé
        function generateGrayShade(index, total) {
            // Generujeme odstíny od světle šedé (80%) do středně šedé (60%)
            const brightness = Math.floor(80 - (index % 5) * 4);
            return 'hsl(0, 0%, ' + brightness + '%)';
        }

        function goUp() {
            if (pathHistory.length > 0) {
                currentData = pathHistory.pop();
                updateDisplay();
            }
        }

        function formatSize(bytes) {
            const units = ['B', 'KB', 'MB', 'GB', 'TB'];
            let size = bytes;
            let unitIndex = 0;
            while (size >= 1024 && unitIndex < units.length - 1) {
                size /= 1024;
                unitIndex++;
            }
            return size.toFixed(1) + ' ' + units[unitIndex];
        }

        function updateChart(data) {
            const ctx = document.getElementById('pieChart').getContext('2d');
            if (chart) {
                chart.destroy();
            }

            document.getElementById('upButton').disabled = pathHistory.length === 0;

            const labels = data.children.map(item => item.path.split('/').pop());
            const sizes = data.children.map(item => item.size);
            const hasChildren = data.children.map(item => item.children && item.children.length > 0);

            const backgroundColors = data.children.map((item, index) => 
                item.children && item.children.length > 0 
                    ? generateColor(index, data.children.length) 
                    : generateGrayShade(index, data.children.length)
            );
            const hoverColors = data.children.map((item, index) => 
                item.children && item.children.length > 0 
                    ? generateColor(index, data.children.length) 
                    : generateGrayShade(index + 1, data.children.length)
            );

            chart = new Chart(ctx, {
                type: 'pie',
                data: {
                    labels: labels,
                    datasets: [{
                        data: sizes,
                        backgroundColor: backgroundColors,
                        hoverBackgroundColor: hoverColors
                    }]
                },
                options: {
                    responsive: true,
                    plugins: {
                        tooltip: {
                            callbacks: {
                                label: function(context) {
                                    const index = context.dataIndex;
                                    const size = formatSize(context.raw);
                                    if (hasChildren[index]) {
                                        return context.label + ': ' + size;
                                    }
                                    return size;
                                }
                            }
                        }
                    },
                    onClick: function(evt, elements) {
                        if (elements.length > 0) {
                            const index = elements[0].index;
                            const clickedItem = data.children[index];
                            if (clickedItem.children && clickedItem.children.length > 0) {
                                pathHistory.push(currentData);
                                currentData = clickedItem;
                                updateDisplay();
                            } else {
                                // Pro neaktivní položky zkopírujeme cestu do schránky
                                navigator.clipboard.writeText(clickedItem.path);
                                // Přidáme vizuální feedback
                                const tooltip = document.createElement('div');
                                tooltip.textContent = 'Cesta zkopírována!';
                                tooltip.style.position = 'fixed';
                                tooltip.style.left = (evt.clientX + 10) + 'px';
                                tooltip.style.top = (evt.clientY - 10) + 'px';
                                tooltip.style.backgroundColor = 'rgba(0,0,0,0.8)';
                                tooltip.style.color = 'white';
                                tooltip.style.padding = '5px 10px';
                                tooltip.style.borderRadius = '4px';
                                tooltip.style.fontSize = '14px';
                                tooltip.style.zIndex = '1000';
                                document.body.appendChild(tooltip);
                                setTimeout(() => {
                                    tooltip.style.opacity = '0';
                                    tooltip.style.transition = 'opacity 0.5s';
                                    setTimeout(() => document.body.removeChild(tooltip), 500);
                                }, 1000);
                            }
                        }
                    }
                }
            });
        }

        function updateDisplay() {
            document.getElementById('currentPath').textContent = currentData.path;
            document.getElementById('currentSize').textContent = formatSize(currentData.size);
            updateChart(currentData);
        }

        function copyPath() {
            navigator.clipboard.writeText(currentData.path);
        }

        updateDisplay();
    </script>
</body>
</html>`
	return html
}
