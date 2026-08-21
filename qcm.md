# QCM final — Développement GO

**Durée :** 1h
**Format :** 20 questions à choix multiples. Cocher la (ou les) bonne(s) réponse(s).
**Notation :** 1 point par question, **0 si une seule mauvaise réponse cochée**.

> Cours couverts : modules 1 à 4 du syllabus (Compilé vs interprété, Langage Go, Bases du langage, Application CLI).

---

## 1. Langages compilés vs interprétés

Quelles affirmations sont correctes ?

- [X] A. Un langage compilé traduit le code source en code machine **avant** l'exécution
- [ ] B. Python distribue typiquement ses programmes sous forme de binaires natifs autonomes
- [ ] C. Un programme Go compilé nécessite une machine virtuelle Go sur la machine cible
- [X] D. Java utilise un modèle intermédiaire : bytecode exécuté par une machine virtuelle (JVM)

## 2. Le langage Go

Quelles affirmations sont correctes ?

- [X] A. Go a été créé chez Google en 2009 (Pike, Thompson, Griesemer)
- [X] B. Docker, Kubernetes et Terraform sont écrits en Go
- [ ] C. Go est un langage à typage dynamique, comme Python
- [X] D. `go build` produit par défaut un binaire statique unique, sans dépendance externe

## 3. La toolchain Go

Quelles affirmations sont correctes ?

- [X] A. `go run main.go` compile puis exécute le programme sans laisser de binaire dans le répertoire
- [X] B. `go build` produit un exécutable
- [X] C. `go fmt` formate le code selon le standard officiel du langage
- [ ] D. `go run` interprète le code ligne par ligne sans jamais le compiler

## 4. Déclarations de variables

Quelles déclarations sont **valides** en Go, placées à l'intérieur d'une fonction ?

- [X] A. `var x int = 5`
- [X] B. `y := 10`
- [X] C. `var z = "texte"`
- [ ] D. `int n = 3;`

## 5. Zero values

Quelles associations type → valeur zéro sont correctes ?

- [X] A. `int` → `0`
- [ ] B. `string` → `"vide"`
- [X] C. `bool` → `false`
- [X] D. `*int` (pointeur) → `nil`

## 6. Slices : qu'affiche ce programme ?

```go
package main

import "fmt"

func main() {
    a := []int{1, 2, 3}
    b := a[:2]
    b = append(b, 99)
    fmt.Println(a)
}
```

- [ ] A. `[1 2 3]`
- [X] B. `[1 2 99]`
- [ ] C. `[1 2 3 99]`
- [ ] D. Le programme panique à l'exécution

## 7. Maps : qu'affiche ce programme ?

```go
package main

import "fmt"

func main() {
    m := map[string]int{"go": 1}
    v, ok := m["rust"]
    fmt.Println(v, ok)
}
```

- [ ] A. `nil false`
- [X] B. `0 false`
- [ ] C. Le programme panique : clé inexistante
- [ ] D. `0 true`

## 8. Pointeurs

Quelles affirmations sont correctes ?

- [X] A. `&x` retourne l'adresse mémoire de `x`
- [X] B. `*p` déréférence le pointeur `p` (accède à la valeur pointée)
- [ ] C. Go permet l'arithmétique de pointeurs comme en C (`p++`)
- [X] D. Passer un pointeur à une fonction permet à celle-ci de modifier la variable d'origine

## 9. Retours multiples : qu'affiche ce programme ?

```go
package main

import "fmt"

func divmod(a, b int) (int, int) {
    return a / b, a % b
}

func main() {
    q, r := divmod(17, 5)
    fmt.Println(q, r)
}
```

- [X] A. `3 2`
- [ ] B. `3.4 2`
- [ ] C. Erreur de compilation : une fonction ne peut retourner qu'une valeur
- [ ] D. `2 3`

## 10. `for`, `range` et `switch`

Quelles affirmations sont correctes ?

- [X] A. `for` est la seule boucle du langage Go (elle couvre while et la boucle infinie)
- [X] B. Sans `fallthrough`, un `case` de `switch` ne « tombe » pas dans le case suivant
- [ ] C. `range` sur une map garantit l'ordre d'insertion des clés
- [ ] D. `while` est un mot-clé réservé de Go

## 11. Interfaces

Quelles affirmations sont correctes ?

- [X] A. Un type implémente une interface **implicitement**, dès qu'il possède toutes ses méthodes
- [ ] B. Il faut le mot-clé `implements` pour lier un type à une interface
- [X] C. `error` est une interface de la stdlib (une seule méthode : `Error() string`)
- [X] D. `any` est un alias de l'interface vide `interface{}`

## 12. Gestion des erreurs

Quelles affirmations sont correctes ?

- [X] A. `fmt.Errorf("... : %w", err)` enveloppe (wrap) l'erreur `err`
- [X] B. `errors.Is(err, cible)` compare `err` à `cible` en remontant la chaîne de wrapping
- [ ] C. La gestion d'erreurs standard en Go passe par un mécanisme `try` / `catch`
- [X] D. Les instructions `defer` s'exécutent en ordre inverse de leur déclaration (LIFO)

## 13. Packages et visibilité

Quelles affirmations sont correctes ?

- [X] A. Un identifiant commençant par une **majuscule** est exporté (visible hors du package)
- [X] B. `go mod init monprojet` crée le fichier `go.mod`
- [X] C. `go get` télécharge et ajoute un package externe aux dépendances du module
- [ ] D. Il faut le mot-clé `public` pour exporter une fonction

## 14. Channels : qu'affiche ce programme ?

```go
package main

import "fmt"

func main() {
    ch := make(chan string)
    go func() { ch <- "pong" }()
    fmt.Println(<-ch)
}
```

- [ ] A. Rien : le programme se termine avant la goroutine
- [X] B. `pong`
- [ ] C. Deadlock : tout envoi sur channel non bufferisé bloque définitivement
- [ ] D. Erreur de compilation

## 15. Goroutines

Quelles affirmations sont correctes ?

- [X] A. Quand `main` se termine, le programme s'arrête même si des goroutines sont encore en cours
- [ ] B. `go f()` attend la fin de `f()` avant de passer à l'instruction suivante
- [X] C. `sync.WaitGroup` permet d'attendre la fin d'un groupe de goroutines
- [ ] D. Une goroutine consomme autant de mémoire qu'un thread OS

## 16. Channels et `select`

Quelles affirmations sont correctes ?

- [X] A. Envoyer sur un channel **non bufferisé** bloque tant qu'aucune goroutine ne lit
- [X] B. `close(ch)` permet à un `for v := range ch` de se terminer proprement
- [X] C. `select` permet d'attendre sur plusieurs opérations de channels à la fois
- [ ] D. **Lire** sur un channel fermé provoque toujours une panique

## 17. Manipulation de fichiers

Quelles affirmations sont correctes ?

- [X] A. `os.ReadFile(path)` lit tout le contenu d'un fichier en mémoire
- [X] B. `os.WriteFile(path, data, 0644)` crée le fichier avec les permissions indiquées (notation octale)
- [X] C. `io.Copy(dst, src)` copie un `io.Reader` vers un `io.Writer` sans tout charger en mémoire
- [ ] D. `os.Remove` supprime récursivement un répertoire non vide

## 18. Réseau HTTP

Quelles affirmations sont correctes ?

- [X] A. `http.Get(url)` retourne `(*http.Response, error)`
- [X] B. Il faut fermer `resp.Body` (par exemple avec `defer resp.Body.Close()`) après usage
- [X] C. `http.ListenAndServe(":8080", nil)` démarre un serveur utilisant le mux par défaut
- [ ] D. Une réponse HTTP 404 fait retourner une `error` non nulle par `http.Get`

## 19. Package `flag` : qu'affiche ce programme, exécuté **sans argument** ?

```go
package main

import (
    "flag"
    "fmt"
)

func main() {
    name := flag.String("name", "monde", "nom à saluer")
    flag.Parse()
    fmt.Println("Bonjour", *name)
}
```

- [X] A. `Bonjour monde`
- [ ] B. `Bonjour` suivi d'une adresse mémoire (`name` est un pointeur)
- [ ] C. Erreur : le flag `-name` est obligatoire
- [ ] D. `Bonjour name`

## 20. Cross-compilation

Quelles affirmations sont correctes ?

- [X] A. `GOOS=windows GOARCH=amd64 go build` produit un `.exe` Windows depuis Linux ou macOS
- [ ] B. La cross-compilation d'un programme Go pur nécessite d'installer un compilateur C croisé
- [X] C. `GOOS=linux GOARCH=arm64` cible Linux sur processeur ARM 64 bits (ex : Raspberry Pi récent)
- [ ] D. Le binaire produit nécessite que Go soit installé sur la machine cible

---

**Fin du QCM.**
Rendre la copie au formateur. Le corrigé est dans `qcm-corrige.md`.
