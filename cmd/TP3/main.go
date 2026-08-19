package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"

	"tp/format"
	"tp/inventory"
)

func donnees() []inventory.File {
	fichiers := []inventory.File{
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
	fichiers = append(fichiers, inventory.File{Name: "notes-perso.txt", Size: 256, Extension: ".txt", Modified: time.Now()})
	inventory.MarquerLesLourds(fichiers, 1_000_000)
	return fichiers
}

func construireRapport(fichiers []inventory.File, sortie string, ext string) (string, error) {
	resultat := fichiers

	if ext != "" {
		filtres, err := inventory.FiltrerParExtension(fichiers, ext)
		if err != nil {
			return "", fmt.Errorf("construction du rapport : %w", err)
		}
		resultat = filtres
	}

	var fmtr format.Formatter
	switch sortie {
	case "texte", "":
		fmtr = format.TextFormatter{Colored: true}
	case "json":
		fmtr = format.JSONFormatter{}
	default:
		return "", fmt.Errorf("construction du rapport : sortie %q inconnue (attendu texte|json)", sortie)
	}

	rapport, err := fmtr.Format(resultat)
	if err != nil {
		return "", fmt.Errorf("construction du rapport : %w", err)
	}
	return rapport, nil
}

func usage() {
	fmt.Println("Usage : gopack [texte|json] [extension]")
	fmt.Println("Exemples : gopack texte        gopack json .pdf")
}

func main() {
	sortie := "texte"
	ext := ""

	if len(os.Args) >= 2 {
		sortie = os.Args[1]
	}
	if len(os.Args) >= 3 {
		ext = os.Args[2]
	}

	fichiers := donnees()

	rapport, err := construireRapport(fichiers, sortie, ext)
	if err != nil {
		switch {
		case errors.Is(err, inventory.ErrAucunResultat):
			color.New(color.FgYellow).Println("Avertissement :", err)
			os.Exit(0)
		case errors.Is(err, inventory.ErrExtensionInvalide):
			color.New(color.FgRed).Fprintln(os.Stderr, "Erreur :", err)
			usage()
			os.Exit(1)
		default:
			color.New(color.FgRed).Fprintln(os.Stderr, "Erreur inattendue :", err)
			os.Exit(1)
		}
	}

	fmt.Println(rapport)

	stats := inventory.StatistiquesParExtension(fichiers)
	if s, ok := stats[".zip"]; ok {
		fmt.Println("Stats .zip :", s)
	}
}
