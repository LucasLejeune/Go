// Package backup copie récursivement un répertoire en préservant
// les permissions des fichiers et des répertoires.
package backup

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Copier copie src vers dst récursivement et retourne le nombre de
// fichiers copiés.
func Copier(src, dst string) (int, error) {
	srcAbs, err := filepath.Abs(src)
	if err != nil {
		return 0, fmt.Errorf("résolution de %s : %w", src, err)
	}
	dstAbs, err := filepath.Abs(dst)
	if err != nil {
		return 0, fmt.Errorf("résolution de %s : %w", dst, err)
	}

	// Piège classique : si dst est à l'intérieur de src, WalkDir copierait
	// la copie dans elle-même à l'infini. On le refuse explicitement.
	if dstAbs == srcAbs || strings.HasPrefix(dstAbs+string(filepath.Separator), srcAbs+string(filepath.Separator)) {
		return 0, fmt.Errorf("la destination %q ne peut pas être à l'intérieur de la source %q", dst, src)
	}

	count := 0
	err = filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)

		info, err := d.Info()
		if err != nil {
			return err
		}

		if d.IsDir() {
			return os.MkdirAll(target, info.Mode().Perm())
		}
		if !d.Type().IsRegular() {
			return nil
		}

		in, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("ouverture %s : %w", path, err)
		}
		defer in.Close()

		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode().Perm())
		if err != nil {
			return fmt.Errorf("création %s : %w", target, err)
		}
		defer out.Close()

		if _, err := io.Copy(out, in); err != nil {
			return fmt.Errorf("copie %s : %w", path, err)
		}

		// Le umask module les permissions à la création : on les réapplique
		// explicitement pour préserver exactement celles de la source.
		if err := os.Chmod(target, info.Mode().Perm()); err != nil {
			return fmt.Errorf("chmod %s : %w", target, err)
		}

		count++
		fmt.Printf("\rCopie... %d fichier(s)", count)
		return nil
	})

	fmt.Println()
	if err != nil {
		return count, fmt.Errorf("copie récursive : %w", err)
	}
	return count, nil
}
