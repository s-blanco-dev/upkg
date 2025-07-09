package pkg

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

type Package struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Arch        string `json:"arch"`
	Filename    string `json:"filename"`
	InstallPath string `json:"install_path"`
}

func (p *Package) Install(repoURL string) error {

	ok, err := IsInstalled(p)
	if err != nil {

		return fmt.Errorf("error verificando si el Mullin se lo tragó: %w", err)
	}
	if ok {
		color.Green("El paquete %s v%s está en el Mullin, salón M103\n", p.Name, p.Version)
		return nil
	}

	// 1. Construir URL del paquete
	url := fmt.Sprintf("%s/%s/%s", repoURL, p.Name, p.Filename)
	fmt.Printf("[#   ] Abriendo puerta del Mullin\n")

	// 2. Descargar el archivo .tar.gz
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error al descargar %s: %w", url, err)
	}
	defer resp.Body.Close()
	fmt.Printf("[##  ] Buscando %s en M203\n", p.Name)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("descarga fallida: %d", resp.StatusCode)
	}

	// 3. Crear carpeta de instalación si no existe
	if err := os.MkdirAll(p.InstallPath, 0755); err != nil {
		return fmt.Errorf("no se pudo crear la carpeta de instalación: %w", err)
	}
	fmt.Printf("[###  ] Succionando %s\n", p.Name)

	// 4. Extraer contenido del .tar.gz
	if err := extractTarGz(resp.Body, p.InstallPath); err != nil {
		return fmt.Errorf("error extrayendo el paquete: %w", err)
	}

	// 5. Registrar el paquete como instalado
	if err := registerInstalledPackage(p); err != nil {
		return fmt.Errorf("no se pudo registrar el paquete: %w", err)
	}
	fmt.Printf("[#### ] Llamando a Jose Abadie\n")

	color.Green("✔ %s instalado correctamente en %s\n", p.Name, p.InstallPath)
	return nil
}

func registerInstalledPackage(pkg *Package) error {
	path := "/etc/upkg/installed.json"

	var installed []*Package

	file, err := os.Open(path)
	if err == nil {
		defer file.Close()
		decodeErr := json.NewDecoder(file).Decode(&installed)
		if decodeErr != nil && decodeErr != io.EOF {
			return decodeErr
		}
	}

	// Agregar el nuevo paquete
	installed = append(installed, pkg)

	// Guardar el archivo nuevamente
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(installed)
}

func extractTarGz(reader io.Reader, targetDir string) error {
	gzr, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tarReader := tar.NewReader(gzr)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // fin del archivo
		}
		if err != nil {
			return err
		}

		path := filepath.Join(targetDir, header.Name)

		switch header.Typeflag {
		case tar.TypeReg:
			outFile, err := os.Create(path)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()

			// Dar permisos ejecutables al archivo
			if err := os.Chmod(path, 0755); err != nil {
				return err
			}
		}
	}

	return nil
}

func IsInstalled(pkg *Package) (bool, error) {
	path := "/etc/upkg/installed.json"

	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	defer file.Close()

	var installed []*Package
	decodeErr := json.NewDecoder(file).Decode(&installed)
	if decodeErr != nil && decodeErr != io.EOF {
		return false, decodeErr
	}

	for _, p := range installed {
		if p.Name == pkg.Name && p.Version == pkg.Version {
			return true, nil
		}
	}

	return false, nil
}
