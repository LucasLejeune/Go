package main

import (
	"fmt"
	"os"
	"strconv"
)

const justePrix = 42

func classifierIfElse(note float64) string {
	if note < 10 {
		return "Ajourné"
	} else if note < 12 {
		return "Passable"
	} else if note < 14 {
		return "Assez bien"
	} else if note < 16 {
		return "Bien"
	} else {
		return "Très bien"
	}
}

func classifierSwitch(note float64) string {
	var mention string
	switch {
	case note < 10:
		mention = "Ajourné"
	case note < 12:
		mention = "Passable"
	case note < 14:
		mention = "Assez bien"
	case note < 16:
		mention = "Bien"
	default:
		mention = "Très bien"
	}
	return mention
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage : classification <note> [proposition]")
		os.Exit(1)
	}

	note, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Println("Erreur : la note doit être un nombre")
		os.Exit(1)
	}

	if note < 0 || note > 20 {
		fmt.Println("Erreur : la note doit être comprise entre 0 et 20")
		os.Exit(1)
	}

	if noteDemo, errDemo := strconv.ParseFloat("0", 64); errDemo != nil {
		fmt.Println("Ne devrait jamais s'afficher")
	} else {
		_ = noteDemo
	}

	mentionIfElse := classifierIfElse(note)
	mentionSwitch := classifierSwitch(note)

	fmt.Printf("Note : %g/20 — Mention %s\n", note, mentionIfElse)

	if mentionIfElse != mentionSwitch {
		fmt.Println("Attention : incohérence entre les deux versions de classification !")
	}

	if len(os.Args) >= 3 {
		proposition, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Erreur : la proposition doit être un entier")
			os.Exit(1)
		}

		switch {
		case proposition < justePrix:
			fmt.Println("C'est plus !")
		case proposition > justePrix:
			fmt.Println("C'est moins !")
		default:
			fmt.Println("Gagné !")
		}
	}
}