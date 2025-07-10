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

	url := fmt.Sprintf("%s/%s/%s", repoURL, p.Name, p.Filename)
	fmt.Printf("[#    ] Abriendo puerta del Mullin\n")

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error al descargar %s: %w", url, err)
	}
	defer resp.Body.Close()
	fmt.Printf("[##   ] Buscando %s en M203\n", p.Name)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("descarga fallida: %d", resp.StatusCode)
	}

	if err := os.MkdirAll(p.InstallPath, 0755); err != nil {
		return fmt.Errorf("no se pudo crear la carpeta de instalación: %w", err)
	}
	fmt.Printf("[###  ] Succionando %s\n", p.Name)

	if err := extractTarGz(resp.Body, p.InstallPath); err != nil {
		return fmt.Errorf("error extrayendo el paquete: %w", err)
	}

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

	// agregar el nuevo paquete
	installed = append(installed, pkg)

	// guardar el archivo luego de agregado el paquete
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
			break
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

			// dar permisos de ejecución al binario
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

func (p *Package) Uninstall() error {
	installed, err := IsInstalled(p)
	if err != nil {
		return fmt.Errorf("error verificando instalación: %w", err)
	}
	if !installed {
		color.Yellow("⚠ El paquete %s v%s no está instalado\n", p.Name, p.Version)
		return nil
	}

	fmt.Printf("[#    ] Buscando %s en el cubículo del Mullin...\n", p.Name)

	// NO SE OLVIDEN DE AGREGAR EL NOMBRE DEL PAQUETE AL InstallPath.
	// Accidentalmente borré todo el directorio /bin
	if err := removePackageFiles(p.InstallPath + "/" + p.Name); err != nil {
		return fmt.Errorf("error eliminando archivos: %w", err)
	}
	fmt.Printf("[##   ] Limpiando restos de %s\n", p.Name)

	if err := unregisterPackage(p); err != nil {
		return fmt.Errorf("error eliminando registro: %w", err)
	}
	fmt.Printf("[###  ] Dando de baja de %s\n", p.Name)

	color.Green("✔ %s v%s desinstalado correctamente\n", p.Name, p.Version)
	return nil
}

func removePackageFiles(installPath string) error {
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		return nil // Si la ruta no existe, no eliminar
	}

	return os.RemoveAll(installPath)
}

// Elimina el paquete de installed.json
func unregisterPackage(pkg *Package) error {
	path := "/etc/upkg/installed.json"

	// Leer lista actual
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil // Si el archivo no existe, no hay nada que hacer
	}
	if err != nil {
		return err
	}
	defer file.Close()

	var installed []*Package
	if err := json.NewDecoder(file).Decode(&installed); err != nil {
		return err
	}

	// Filtrar el paquete a eliminar
	filtered := make([]*Package, 0)
	for _, p := range installed {
		if p.Name != pkg.Name || p.Version != pkg.Version {
			filtered = append(filtered, p)
		}
	}

	// Guardar lista actualizada
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(filtered)
}
