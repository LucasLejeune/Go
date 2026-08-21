// Package scan inventorie récursivement un répertoire et calcule des
// statistiques par extension.
package scan

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
)

type Stats struct {
	Nombre       int
	TailleTotale int64
}

type Rapport struct {
	ParExtension  map[string]Stats
	TotalFichiers int
	TotalOctets   int64
}

// Scanner parcourt dir récursivement et retourne le rapport agrégé.
// Le calcul est séparé de l'affichage pour rester testable.
func Scanner(dir string) (Rapport, error) {
	rapport := Rapport{ParExtension: make(map[string]Stats)}

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.Type().IsRegular() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		ext := filepath.Ext(path)
		if ext == "" {
			ext = "(sans extension)"
		}

		s := rapport.ParExtension[ext]
		s.Nombre++
		s.TailleTotale += info.Size()
		rapport.ParExtension[ext] = s

		rapport.TotalFichiers++
		rapport.TotalOctets += info.Size()
		return nil
	})
	if err != nil {
		return Rapport{}, fmt.Errorf("parcours de %s : %w", dir, err)
	}
	return rapport, nil
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

// Afficher imprime le rapport trié par taille décroissante avec un total général.
func Afficher(r Rapport) {
	type ligne struct {
		ext   string
		stats Stats
	}

	var lignes []ligne
	for ext, s := range r.ParExtension {
		lignes = append(lignes, ligne{ext, s})
	}
	sort.Slice(lignes, func(i, j int) bool {
		return lignes[i].stats.TailleTotale > lignes[j].stats.TailleTotale
	})

	fmt.Printf("%-20s %-10s %s\n", "EXTENSION", "FICHIERS", "TAILLE")
	for _, l := range lignes {
		fmt.Printf("%-20s %-10d %s\n", l.ext, l.stats.Nombre, tailleLisible(l.stats.TailleTotale))
	}
	fmt.Printf("\nTotal : %d fichier(s), %s\n", r.TotalFichiers, tailleLisible(r.TotalOctets))
}
