package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type Result struct {
	Name string
	Err  error
}

func fetchList(serverURL string) ([]FileInfo, error) {
	resp, err := http.Get(serverURL + "/files")
	if err != nil {
		return nil, fmt.Errorf("récupération de la liste : %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("récupération de la liste : statut %d", resp.StatusCode)
	}

	var files []FileInfo
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, fmt.Errorf("décodage JSON : %w", err)
	}
	return files, nil
}

func downloadFile(serverURL, outDir, name string) error {
	resp, err := http.Get(serverURL + "/files/" + name)
	if err != nil {
		return fmt.Errorf("requête %s : %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("téléchargement %s : statut %d", name, resp.StatusCode)
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("création du dossier %s : %w", outDir, err)
	}

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

func main() {
	serverURL := flag.String("server", "http://localhost:8080", "URL du serveur de fichiers")
	outDir := flag.String("out", "downloads", "répertoire de destination")
	n := flag.Int("n", 4, "nombre de téléchargements simultanés (sémaphore)")
	flag.Parse()

	files, err := fetchList(*serverURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur :", err)
		os.Exit(1)
	}

	total := len(files)
	results := make(chan Result)
	sem := make(chan struct{}, *n)

	var wg sync.WaitGroup
	for _, f := range files {
		wg.Add(1)
		go func(f FileInfo) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			err := downloadFile(*serverURL, *outDir, f.Name)
			results <- Result{Name: f.Name, Err: err}
		}(f)
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
	if echecs > 0 {
		os.Exit(1)
	}
}
