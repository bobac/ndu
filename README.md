# NDU - Command line utility for disk usage analysis

A utility for analyzing directory sizes on disk, similar to the Unix `du` command.

## Features
- Displays directory sizes in human-readable format (e.g. 1.2 GB)
- Shows only the largest directories
- Recursive analysis with configurable depth
- JSON export of directory structure and sizes
- Cross-platform support (Windows, Linux, macOS)

## Usage
ndu [switches] [directory]

### Switches
- `-h` - Displays directory sizes in human-readable format (e.g. 24 KB, 2.2 GB)
- `-n count` - Shows only "count" largest directories at each level
- `-r depth` - Performs recursive analysis of directories up to "depth" level
- `-d count` - Limits the number of directories for recursive analysis to "count" largest ones
- `-v` - Shows currently processed directory
- `-j file.json` - Exports results to JSON file (sizes in bytes)
- `-help` - Shows help message

## Examples
```bash
# Show top 3 largest directories in human-readable format
ndu -h -n 3 /

# Recursive analysis up to depth 2, showing top 3 directories at each level
ndu -h -n 3 -r 2 /

# Recursive analysis up to depth 2, showing top 3 directories but only analyzing the largest one
ndu -h -n 3 -r 2 -d 1 /

# Export directory analysis to JSON file
ndu -j results.json /

# Combine JSON export with other options
ndu -h -n 3 -r 2 -j results.json /
```

## JSON Output Format
The JSON output contains a hierarchical structure of directories with their sizes in bytes:
```json
{
  "path": "/path/to/directory",
  "size": 1234567890,
  "children": [
    {
      "path": "/path/to/directory/subdir1",
      "size": 123456789,
      "children": []
    }
  ]
}
```

## Installation
```bash
go install github.com/bobac/ndu/cmd/ndu@latest
```

## License
MIT License - see [LICENSE](LICENSE) file for details 