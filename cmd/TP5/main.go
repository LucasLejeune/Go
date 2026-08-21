package main

import (
	"flag"
	"fmt"
	"os"

	"tp/internal/backup"
	"tp/internal/fetch"
	"tp/internal/scan"
)

// version est injectée au build avec -ldflags "-X main.version=..."
var version = "dev"

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(2)
	}

	var err error

	switch os.Args[1] {
	case "scan":
		err = runScan(os.Args[2:])
	case "backup":
		err = runBackup(os.Args[2:])
	case "fetch":
		err = runFetch(os.Args[2:])
	case "version":
		fmt.Println("gopack", version)
		return
	case "help", "-h", "--help":
		printHelp()
		return
	default:
		fmt.Fprintf(os.Stderr, "gopack : commande inconnue %q\n\n", os.Args[1])
		printHelp()
		os.Exit(2)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "gopack :", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`gopack - outil CLI de sauvegarde et transfert de fichiers

Usage :
  gopack <commande> [arguments]

Commandes :
  scan <dir>                     Inventaire récursif d'un répertoire
  backup <src> <dst>             Copie récursive avec permissions préservées
  fetch [-o dir] [-n N] <url>... Téléchargement concurrent de fichiers
  version                        Affiche la version
  help                           Affiche cette aide

Utilisez 'gopack <commande> -h' pour l'aide d'une sous-commande.`)
}

func runScan(args []string) error {
	fs := flag.NewFlagSet("scan", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage : gopack scan <dir>")
		fmt.Fprintln(os.Stderr, "\nInventorie récursivement un répertoire et affiche les statistiques par extension.")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		fs.Usage()
		os.Exit(2)
	}

	rapport, err := scan.Scanner(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("scan : %w", err)
	}
	scan.Afficher(rapport)
	return nil
}

func runBackup(args []string) error {
	fs := flag.NewFlagSet("backup", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage : gopack backup <src> <dst>")
		fmt.Fprintln(os.Stderr, "\nCopie récursivement <src> vers <dst> en préservant les permissions.")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 2 {
		fs.Usage()
		os.Exit(2)
	}

	n, err := backup.Copier(fs.Arg(0), fs.Arg(1))
	if err != nil {
		return fmt.Errorf("backup : %w", err)
	}
	fmt.Printf("Sauvegarde terminée : %d fichier(s) copié(s)\n", n)
	return nil
}

func runFetch(args []string) error {
	fs := flag.NewFlagSet("fetch", flag.ExitOnError)
	outDir := fs.String("o", "downloads", "répertoire de destination")
	workers := fs.Int("n", 4, "nombre de téléchargements simultanés")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage : gopack fetch [-o dir] [-n N] <url>...")
		fmt.Fprintln(os.Stderr, "\nTélécharge une ou plusieurs URLs en parallèle (sémaphore -n).")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		fs.Usage()
		os.Exit(2)
	}

	echecs, err := fetch.Telecharger(fs.Args(), *outDir, *workers)
	if err != nil {
		return fmt.Errorf("fetch : %w", err)
	}
	if echecs > 0 {
		os.Exit(1)
	}
	return nil
}
