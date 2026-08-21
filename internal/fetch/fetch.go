// Package fetch télécharge une liste d'URLs en parallèle, avec une
// concurrence limitée par sémaphore. Réutilise l'architecture du TP4.
package fetch

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type Result struct {
	URL  string
	Name string
	Err  error
}

func nomDepuisURL(url string) string {
	return filepath.Base(url)
}

func telechargerUn(url, outDir string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("requête : %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("statut %d", resp.StatusCode)
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("création du dossier %s : %w", outDir, err)
	}

	name := nomDepuisURL(url)
	out, err := os.OpenFile(
		filepath.Join(outDir, name),
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0o644,
	)
	if err != nil {
		return fmt.Errorf("création du fichier %s : %w", name, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("écriture %s : %w", name, err)
	}
	return nil
}

// Telecharger télécharge toutes les urls en parallèle (concurrence limitée
// à workers téléchargements simultanés) et retourne le nombre d'échecs.
// Une URL en échec n'arrête jamais le lot.
func Telecharger(urls []string, outDir string, workers int) (int, error) {
	total := len(urls)
	results := make(chan Result)
	sem := make(chan struct{}, workers)

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			err := telechargerUn(url, outDir)
			results <- Result{URL: url, Name: nomDepuisURL(url), Err: err}
		}(url)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	succes, echecs, i := 0, 0, 0
	for res := range results {
		i++
		if res.Err != nil {
			echecs++
			fmt.Fprintf(os.Stderr, "[%d/%d] %s : %v\n", i, total, res.Name, res.Err)
			continue
		}
		succes++
		fmt.Printf("[%d/%d] %s : ok\n", i, total, res.Name)
	}

	fmt.Printf("\nBilan : %d succès, %d échecs\n", succes, echecs)
	return echecs, nil
}
