package format

import "tp/inventory"

type Formatter interface {
	Format(fichiers []inventory.File) (string, error)
}
