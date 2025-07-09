package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ucunix/upkg/pkg"

	"github.com/fatih/color"
)

// Literalmente lo copié todo de Stack Overflow
func LoadFromURL(url string) (*pkg.Package, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("el Mullin se te cayó encima: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("busqué en cada rincón del Mullin y no encontré nada: %d", resp.StatusCode)
	}

	var pkg pkg.Package
	if err := json.NewDecoder(resp.Body).Decode(&pkg); err != nil {
		return nil, fmt.Errorf("error al abrir puerta del Mullin: %w", err)
	}

	return &pkg, nil
}

func ListRemotePackages(indexURL string) error {
	resp, err := http.Get(indexURL)
	if err != nil {
		return fmt.Errorf("No se pudo acceder a la cartelera del Mullin: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("el Mullin respondió con código %d", resp.StatusCode)
	}

	var packages []*pkg.Package
	if err := json.NewDecoder(resp.Body).Decode(&packages); err != nil {
		return fmt.Errorf("error al decodificar index.json: %w", err)
	}

	if len(packages) == 0 {
		fmt.Println("No hay paquetes disponibles en el repositorio.")
		return nil
	}

	fmt.Println("📦 Paquetes disponibles en repositorio:")
	for _, p := range packages {
		color.HiGreen("  - %s (%s) → %s\n", p.Name, p.Version, p.Filename)
	}

	return nil
}
