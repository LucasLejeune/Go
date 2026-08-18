# TP1 — Environnement Go & premiers programmes

Repo contenant les 4 exécutables du TP1, chacun dans son propre sous-dossier de `cmd/` (convention Go standard pour un module multi-binaires).

## Arborescence

```
tp1-go/
├── go.mod
├── .gitignore
├── README.md
└── cmd/
    ├── hello/main.go
    ├── convertisseur/main.go
    ├── classification/main.go
    └── calculatrice/main.go
```

## Build & run

```bash
# Etape 1
go run ./cmd/hello
go build -o bin/hello ./cmd/hello

# Etape 2
go run ./cmd/convertisseur 42

# Etape 3
go run ./cmd/classification 15
go run ./cmd/classification 15 40   # bonus juste prix

# Etape 4 (livrable)
go build -o bin/calc ./cmd/calculatrice
./bin/calc 5 + 3
./bin/calc 10 / 4
./bin/calc 1 / 0
./bin/calc version
```

Compiler tous les exécutables d'un coup depuis la racine :

```bash
go build ./...
```

## Réponses aux questions checkpoints

- **Taille du binaire `hello`** : plusieurs Mo alors que le code source fait 5 lignes, car Go produit un binaire **statique** : le runtime Go (goroutines, garbage collector, gestion mémoire) et les paquets de la bibliothèque standard utilisés (`fmt`) sont embarqués directement dans l'exécutable. Il n'y a aucune dépendance dynamique à installer sur la machine cible.

- **Pourquoi `strconv.ParseFloat` est nécessaire** : `os.Args` retourne uniquement des `string`. Go est un langage à typage fort qui ne fait **jamais** de conversion implicite entre `string` et types numériques ; il faut donc explicitement parser la chaîne avec `strconv`.

- **Pourquoi le `switch` Go n'a pas besoin de `break`** : contrairement au C ou à JavaScript, chaque `case` d'un `switch` Go s'arrête automatiquement après exécution (pas de fall-through implicite). Pour forcer le comportement inverse, il faut utiliser explicitement le mot-clé `fallthrough`.

- **Scope du `if` à instruction courte** (`if x := ...; cond { }`) : la variable déclarée (`x`) n'existe que dans le bloc `if` et ses éventuels `else if`/`else` associés. Elle n'est pas visible en dehors de cette structure conditionnelle.

## Pièges rencontrés

- `int / int` fait une division entière en Go : la conversion `classification.go`/`convertisseur.go` reste en `float64` de bout en bout pour éviter les troncatures indésirables.
- `./calc 5 * 3` sans guillemets échoue car le shell interprète `*` comme un glob de fichiers ; d'où l'opérateur `x` en principal et `*` en alias géré dans le `switch`.
