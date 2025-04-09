# NDU - Command line utilita pro zjištění obsazenosti disku

Utilita pro analýzu velikosti adresářů na disku, podobná unixovému příkazu `du`.

## Instalace

```bash
go install github.com/bobac/ndu/cmd/ndu@latest
```

## Použití

```bash
ndu [switche] [adresář]
```

### Switche
- `-h` - Vypisuje velikosti adresářů v lidsky čitelné formě (např. 24 KB, 2.2 GB)
- `-n počet` - Vypíše jen "počet" největších adresářů na každé úrovni
- `-r hloubka` - Provede rekurzivní analýzu adresářů do úrovně "hloubka"
- `-d počet` - Omezí počet adresářů rekurzivní analýzy nejvyšší počet "počet" největších

### Příklady

```bash
# Zobrazení 10 největších adresářů v aktuálním adresáři
ndu -h -n 10

# Rekurzivní analýza kořenového adresáře
ndu -h -n 3 -r -d 1 /
```

## Build

Pro build na všech podporovaných platformách použijte:

### Linux/MacOS
```bash
chmod +x build.sh
./build.sh
```

### Windows
```powershell
.\build.ps1
```

Výsledné binární soubory budou umístěny v adresáři `bin/`.

## Podporované platformy
- Windows
- Linux (x86_64, ARM64)
- MacOS (Intel, Apple Silicon)

## Licence
MIT 