package format

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"

	"tp/inventory"
)

var ErrRienAFormater = errors.New("aucun fichier à formater")

type TextFormatter struct {
	Colored bool
}

func (t TextFormatter) Format(fichiers []inventory.File) (string, error) {
	if len(fichiers) == 0 {
		return "", ErrRienAFormater
	}

	var b strings.Builder

	entete := fmt.Sprintf("%-25s %-12s %-12s %s", "NOM", "TAILLE", "MODIFIÉ", "TAG")
	if t.Colored {
		entete = color.New(color.FgCyan, color.Bold).Sprint(entete)
	}
	b.WriteString(entete)
	b.WriteString("\n")

	for _, f := range fichiers {
		tag := f.Tag
		if t.Colored && tag != "" {
			tag = color.New(color.FgYellow).Sprint(tag)
		}
		fmt.Fprintf(&b, "%-25s %-12s %-12s %s\n", f.Name, f.TailleLisible(), f.Modified.Format("2006-01-02"), tag)
	}

	return b.String(), nil
}
