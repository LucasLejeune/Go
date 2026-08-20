package inventory

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	ErrExtensionInvalide = errors.New("extension invalide (attendu : \".ext\")")
	ErrTailleNegative    = errors.New("taille minimale négative")
	ErrAucunResultat     = errors.New("aucun fichier ne correspond")
)

type File struct {
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	Extension string    `json:"extension"`
	Modified  time.Time `json:"modified"`
	Tag       string    `json:"tag,omitempty"`
}

type Stats struct {
	Nombre       int
	TailleTotale int64
}

func (s Stats) String() string {
	return fmt.Sprintf("%d fichier(s), %s", s.Nombre, tailleLisible(s.TailleTotale))
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

func MarquerLesLourds(fichiers []File, tailleMin int64) int {
	count := 0
	for i := range fichiers {
		if fichiers[i].Size >= tailleMin {
			fichiers[i].Marquer("a-archiver")
			count++
		}
	}
	return count
}

func FiltrerParExtension(fichiers []File, ext string) ([]File, error) {
	if ext == "" || !strings.HasPrefix(ext, ".") {
		return nil, fmt.Errorf("filtre extension %q : %w", ext, ErrExtensionInvalide)
	}

	var resultat []File
	for _, f := range fichiers {
		if f.Extension == ext {
			resultat = append(resultat, f)
		}
	}

	if len(resultat) == 0 {
		return nil, fmt.Errorf("filtre extension %q : %w", ext, ErrAucunResultat)
	}
	return resultat, nil
}

func FiltrerParTailleMin(fichiers []File, tailleMin int64) ([]File, error) {
	if tailleMin < 0 {
		return nil, fmt.Errorf("filtre taille %d : %w", tailleMin, ErrTailleNegative)
	}

	var resultat []File
	for _, f := range fichiers {
		if f.Size >= tailleMin {
			resultat = append(resultat, f)
		}
	}

	if len(resultat) == 0 {
		return nil, fmt.Errorf("filtre taille %d : %w", tailleMin, ErrAucunResultat)
	}
	return resultat, nil
}

func StatistiquesParExtension(fichiers []File) map[string]Stats {
	stats := make(map[string]Stats)
	for _, f := range fichiers {
		s := stats[f.Extension]
		s.Nombre++
		s.TailleTotale += f.Size
		stats[f.Extension] = s
	}
	return stats
}

func TrierParTaille(fichiers []File) []File {
	tries := make([]File, len(fichiers))
	copy(tries, fichiers)
	sort.Slice(tries, func(i, j int) bool {
		return tries[i].Size > tries[j].Size
	})
	return tries
}
