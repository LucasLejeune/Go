package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "version" {
		fmt.Println("calc v1.0")
		return
	}

	if len(os.Args) != 4 {
		fmt.Println("Usage : calc <nombre> <opérateur> <nombre>")
		fmt.Println("Opérateurs : + - x /")
		os.Exit(1)
	}

	operateur := os.Args[2]

	if operateur == "%" {
		a, errA := strconv.Atoi(os.Args[1])
		b, errB := strconv.Atoi(os.Args[3])
		if errA != nil || errB != nil {
			fmt.Println("Erreur : le modulo nécessite deux nombres entiers")
			os.Exit(1)
		}
		if b == 0 {
			fmt.Println("Erreur : modulo par zéro")
			os.Exit(1)
		}
		fmt.Printf("%d %s %d = %d\n", a, operateur, b, a%b)
		return
	}

	a, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Println("Erreur : le premier opérande doit être un nombre")
		os.Exit(1)
	}

	b, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		fmt.Println("Erreur : le second opérande doit être un nombre")
		os.Exit(1)
	}

	var resultat float64

	switch operateur {
	case "+":
		resultat = a + b
	case "-":
		resultat = a - b
	case "x", "*":
		resultat = a * b
	case "/":
		if b == 0 {
			fmt.Println("Erreur : division par zéro")
			os.Exit(1)
		}
		resultat = a / b
	default:
		fmt.Printf("Erreur : opérateur inconnu %q\n", operateur)
		os.Exit(1)
	}

	fmt.Printf("%g %s %g = %g\n", a, operateur, b, resultat)
}