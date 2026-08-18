package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type File struct {
	Name      string
	Size      int64
	Extension string
	Modified  time.Time
	Tag       string
}

type Stats struct {
	Nombre       int
	TailleTotale int64
}

var inventaire = []File{
	{Name: "main.go", Size: 4_096, Extension: ".go", Modified: time.Date(2026, 7, 1, 10, 15, 0, 0, time.UTC)},
	{Name: "utils.go", Size: 2_048, Extension: ".go", Modified: time.Date(2026, 7, 3, 16, 40, 0, 0, time.UTC)},
	{Name: "README.md", Size: 8_192, Extension: ".md", Modified: time.Date(2026, 6, 20, 9, 0, 0, 0, time.UTC)},
	{Name: "lisezmoi.md", Size: 1_024, Extension: ".md", Modified: time.Date(2026, 6, 21, 11, 30, 0, 0, time.UTC)},
	{Name: "rapport.pdf", Size: 1_254_400, Extension: ".pdf", Modified: time.Date(2026, 6, 12, 9, 30, 0, 0, time.UTC)},
	{Name: "presentation.pdf", Size: 3_670_016, Extension: ".pdf", Modified: time.Date(2026, 5, 30, 14, 0, 0, 0, time.UTC)},
	{Name: "logo.png", Size: 46_080, Extension: ".png", Modified: time.Date(2026, 4, 2, 8, 45, 0, 0, time.UTC)},
	{Name: "photo-equipe.png", Size: 2_936_012, Extension: ".png", Modified: time.Date(2026, 4, 2, 8, 50, 0, 0, time.UTC)},
	{Name: "todo.txt", Size: 512, Extension: ".txt", Modified: time.Date(2026, 7, 18, 18, 5, 0, 0, time.UTC)},
	{Name: "backup.zip", Size: 15_728_640, Extension: ".zip", Modified: time.Date(2026, 1, 15, 3, 0, 0, 0, time.UTC)},
	{Name: "archive-2025.zip", Size: 44_040_192, Extension: ".zip", Modified: time.Date(2025, 12, 31, 23, 59, 0, 0, time.UTC)},
	{Name: "config.txt", Size: 730, Extension: ".txt", Modified: time.Date(2026, 2, 10, 12, 0, 0, 0, time.UTC)},
}

func tailleLisible(octets int64) string {
	switch {
	case octets >= 1024*1024:
		return fmt.Sprintf("%.1f Mo", float64(octets)/(1024*1024))
	case octets >= 1024:
		return fmt.Sprintf("%.1f Ko", float64(octets)/1024)
	default:
		return fmt.Sprintf("%d o", octets)
	}
}

func (f File) TailleLisible() string {
	return tailleLisible(f.Size)
}

func (f *File) Renommer(nouveauNom string) {
	f.Name = nouveauNom
	f.Extension = filepath.Ext(nouveauNom)
}

func (f *File) Marquer(tag string) {
	f.Tag = tag
}

func marquerLesLourds(fichiers []File, tailleMin int64) int {
	count := 0
	for i := range fichiers {
		if fichiers[i].Size >= tailleMin {
			fichiers[i].Marquer("a-archiver")
			count++
		}
	}
	return count
}

func afficher(fichiers []File) {
	fmt.Printf("%-25s %-12s %-12s %s\n", "NOM", "TAILLE", "MODIFIÉ", "TAG")
	for _, f := range fichiers {
		fmt.Printf("%-25s %-12s %-12s %s\n", f.Name, f.TailleLisible(), f.Modified.Format("2006-01-02"), f.Tag)
	}
}

func filtrerParExtension(fichiers []File, ext string) ([]File, error) {
	if ext == "" || !strings.HasPrefix(ext, ".") {
		return nil, errors.New("extension invalide : elle doit commencer par un point")
	}
	var resultat []File
	for _, f := range fichiers {
		if f.Extension == ext {
			resultat = append(resultat, f)
		}
	}
	return resultat, nil
}

func filtrerParTailleMin(fichiers []File, tailleMin int64) ([]File, error) {
	if tailleMin < 0 {
		return nil, errors.New("tailleMin ne peut pas être négative")
	}
	var resultat []File
	for _, f := range fichiers {
		if f.Size >= tailleMin {
			resultat = append(resultat, f)
		}
	}
	return resultat, nil
}

func statistiquesParExtension(fichiers []File) map[string]Stats {
	stats := make(map[string]Stats)
	for _, f := range fichiers {
		s := stats[f.Extension]
		s.Nombre++
		s.TailleTotale += f.Size
		stats[f.Extension] = s
	}
	return stats
}

func trierParTaille(fichiers []File) []File {
	tries := make([]File, len(fichiers))
	copy(tries, fichiers)
	sort.Slice(tries, func(i, j int) bool {
		return tries[i].Size > tries[j].Size
	})
	return tries
}

func supprimer(fichiers []File, nom string) ([]File, error) {
	for i, f := range fichiers {
		if f.Name == nom {
			return append(fichiers[:i], fichiers[i+1:]...), nil
		}
	}
	return nil, fmt.Errorf("fichier %q introuvable", nom)
}

func plusRecentQue(fichiers []File, date time.Time) []File {
	var resultat []File
	for _, f := range fichiers {
		if f.Modified.After(date) {
			resultat = append(resultat, f)
		}
	}
	return resultat
}

func main() {
	fmt.Println("TP2 — inventaire de fichiers : à vous de jouer !")

	fmt.Println("\n=== Étape 1 : struct File ===")
	f := File{Name: "test.txt", Size: 100, Extension: ".txt", Modified: time.Now()}
	fmt.Println(f)
	fmt.Printf("%+v\n", f)

	fmt.Println("\n=== Étape 2 : inventaire ===")
	fmt.Println("Nombre de fichiers :", len(inventaire))
	for i, fi := range inventaire {
		fmt.Printf("%d. %s (%d octets)\n", i+1, fi.Name, fi.Size)
	}
	inventaire = append(inventaire, File{Name: "notes-perso.txt", Size: 256, Extension: ".txt", Modified: time.Now()})
	fmt.Println("Nombre de fichiers après append :", len(inventaire))

	fmt.Println("\n=== Étape 3 : filtrage ===")
	goFiles, err := filtrerParExtension(inventaire, ".go")
	if err != nil {
		fmt.Println("Erreur de filtrage :", err)
		os.Exit(1)
	}
	fmt.Println(".go trouvés :", len(goFiles))

	if _, err := filtrerParExtension(inventaire, "go"); err != nil {
		fmt.Println("Erreur attendue (pas de point) :", err)
	}

	gros, err := filtrerParTailleMin(inventaire, 1_000_000)
	if err != nil {
		fmt.Println("Erreur de filtrage :", err)
		os.Exit(1)
	}
	fmt.Println("Fichiers ≥ 1 Mo :", len(gros))

	if _, err := filtrerParTailleMin(inventaire, -1); err != nil {
		fmt.Println("Erreur attendue (taille négative) :", err)
	}

	pngs, _ := filtrerParExtension(inventaire, ".png")
	pngsLourds, _ := filtrerParTailleMin(pngs, 1_000_000)
	fmt.Println(".png de plus de 1 Mo :", len(pngsLourds))
	for _, p := range pngsLourds {
		fmt.Println(" -", p.Name)
	}
	fmt.Println("Inventaire non modifié, taille toujours :", len(inventaire))

	fmt.Println("\n=== Étape 4 : statistiques par extension ===")
	stats := statistiquesParExtension(inventaire)

	var extensions []string
	for ext := range stats {
		extensions = append(extensions, ext)
	}
	sort.Strings(extensions)

	for _, ext := range extensions {
		s := stats[ext]
		fmt.Printf("%-6s %d fichier(s), %s\n", ext, s.Nombre, tailleLisible(s.TailleTotale))
	}

	if s, ok := stats[".zip"]; ok {
		fmt.Println("Les .zip pèsent", s.TailleTotale, "octets")
	}
	if _, ok := stats[".xyz"]; !ok {
		fmt.Println("Aucune extension .xyz dans l'inventaire (zero value sans panique)")
	}

	fmt.Println("\n=== Étape 5 : tri par taille ===")
	tries := trierParTaille(inventaire)
	fmt.Println("Top 3 des fichiers les plus lourds :")
	for i := 0; i < 3 && i < len(tries); i++ {
		fmt.Printf("%d. %s (%s)\n", i+1, tries[i].Name, tries[i].TailleLisible())
	}
	fmt.Println("Premier élément de l'inventaire d'origine (doit être inchangé) :", inventaire[0].Name)

	fmt.Println("\n=== Étape 6 : affichage formaté ===")
	afficher(inventaire)

	fmt.Println("\n=== Étape 7 : méthodes & pointer receiver ===")
	for i := range inventaire {
		if inventaire[i].Name == "todo.txt" {
			inventaire[i].Renommer("todo-2026.md")
		}
	}
	nbMarques := marquerLesLourds(inventaire, 1_000_000)
	fmt.Println("Fichiers marqués 'a-archiver' :", nbMarques)
	afficher(inventaire)
}
