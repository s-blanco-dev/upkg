package local

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"ucunix/upkg/pkg"

	"github.com/fatih/color"
)

func ListInstalledPackages(installedPath string) error {
	file, err := os.Open(installedPath)
	if os.IsNotExist(err) {
		color.Yellow("No hay paquetes en tu sistema. -> El Mullin está vacío")
		return nil
	}
	if err != nil {
		return fmt.Errorf("error al abrir la base de paquetes: %w", err)
	}
	defer file.Close()

	var installed []*pkg.Package
	decodeErr := json.NewDecoder(file).Decode(&installed)
	if decodeErr != nil {
		return fmt.Errorf("error al leer los paquetes instalados: %w", decodeErr)
	}

	if len(installed) == 0 {
		color.Yellow("No hay paquetes instalados.")
		return nil
	}

	fmt.Println("📦 Paquetes instalados:")
	for _, p := range installed {
		color.Green("  - %s (%s) → %s\n", p.Name, p.Version, p.InstallPath)
	}

	return nil
}

func GetInstalledPkg(name string) (*pkg.Package, error) {
	path := "/etc/upkg/installed.json"

	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var installed []*pkg.Package
	decodeErr := json.NewDecoder(file).Decode(&installed)
	if decodeErr != nil && decodeErr != io.EOF {
		return nil, decodeErr
	}

	for _, p := range installed {
		if p.Name == name {
			return p, nil
		}
	}

	return nil, nil
}
