package local

import (
	"encoding/json"
	"fmt"
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
