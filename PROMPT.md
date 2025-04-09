# NDU - Command line utilita pro zjištění obsazenosti disku

## Funkce
Utilita zjišťuje, které adresáře na disku zabírají nejvíc místa a vypisuje ty největší v pořadí od největšího.
Odpovídá funkci unixového příkazu: `du -h --max-depth=1 / | sort -hr | head -n 10`. Tomuto příkazu by odpovídalo následující vloání ndu: `ndu -h -n 10 /` na unix-like systémech, nebo `ndu -h -n 10 "c:\"` na windows.
V případě použití přepínaře `-r` vstoupí rekurzívně do každého adresáře (až do počtu určeného přepínačem `-n`, což je maximum) nebo dále omezí tento počet pomocí přepínače `-d`. V takovém případě vždy vypisuje v daném adresář `n` velikostí podadresářů, ale následně rekurzivně vstoupí jen do prvních `d` adresářů.

## Příklady výstupů

Vstup: `ndu -h -n 3 /`

Výstup:
```
/var    234 GB
/home   123 GB
/bin     26 MB
```

Vstup: `ndu -h -n 3 -r 1 -d 1 /`

Výstup:
```
var/    234 GB
home/   123 GB
bin/     26 MB

=> /var
var/lib/    189 GB
var/log/     21 GB
var/tmp/    143 MB
```

## Použití
ndu [switche] [adresář]

### Switche
`-h` Vypisuje velikosti adresářů v lidsky čitelné formě, tedy např.: 24 KB, 2.2 GB, 26 B apod. Pokud není uvedeno, vypisuje jen čísla
`-n počet` Vypíše jen "počet" největších adresářů na každé úrovni. Pokud není uvedeno, vypíše vše.
`-r hloubka` Provede sám sebe rekurzivně na všechny adresáře do hloubky "hloubka". Může být omezeno pomocí `-d počet`
`-d počet` Pro "počet" největších adresářů v každé úrovni provede rekurzivně sám sebe na daný adresář se stejnými parametry. Pokud není uvedeno projde všechny. Vyžaduje `-r`

### Adresář
Pokud není parametr adresář uveden začínáme počítat velikost od aktuálního adresáře.

## Struktura projektu
Připrav celý projekt v jazyce GO jako modul, který bude umístěn na github.com/bobac/ndu, aby v `./cmd/ndu/main.go` byl samotný spustitelnů kód, ale všechny funkce byly volány z tohoto modulu. Chci mít možnost celou funkcionalitu přidávat i do jiných projektů.

## Platformy
- Windows
- Linux
- Linux ARM
- MacOS - Intel
- MacOS - Apple

Potřebuji zajistit, aby kód fungoval na všech výše uvedených platformách, pokud je třeba rozdílné API, a taky připravit 2 build skripty v Powershell a bash, který do adresáře ./bin připraví spustitelné soubory pro všechny platformy.

## GIT
Připrav soubor `.gitignore` tak, aby se vyhnul artefaktům, které nechceme dávat do repozitáře, např: spustitelné soubory.
