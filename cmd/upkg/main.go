package main

import (
	"fmt"
	"log"
	"os"
	"ucunix/upkg/local"
	"ucunix/upkg/utils"

	"github.com/fatih/color"
)

const repoUrl = "https://s-blanco-dev.github.io/ucunix_pkg"
const installedDatabase = "/etc/upkg/installed.json"

func main() {

	switch os.Args[1] {
	case "install":
		installPackage(os.Args[2])
	case "list":
		listPackages()
	case "installed":
		listInstalled()
	case "help":
		printHelp()
	default:
		color.Red("¿Tenés la cabeza plana?\n")
		printHelp()
	}
}

func installPackage(pkgName string) {
	pak, err := utils.LoadFromURL(repoUrl + "/" + pkgName + "/PKG.json")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if err := pak.Install(repoUrl); err != nil {
		log.Fatalf("Error instalando el paquete: %v", err)
	}
}

func listPackages() {
	if err := utils.ListRemotePackages(repoUrl); err != nil {
		log.Fatalf("Error listando paquetes desde repositorio externo: %v", err)
	}
}

func listInstalled() {
	if err := local.ListInstalledPackages(installedDatabase); err != nil {
		log.Fatalf("Error listando paquetes desde repositorio externo: %v", err)
	}
}

func printHelp() {
	color.Cyan("UCU package manager (upkg) -> v.1.1\n")
	fmt.Println("Santiago Blanco 2025")
	fmt.Println("---------------------")
	fmt.Println("Opciones:")
	fmt.Println("Instalar paquete desde Mullin ->")
	color.Green("upkg install <paquete>")

	fmt.Println("Listar paquetes disponibles en el Mullin ->")
	color.Green("upkg list")

	fmt.Println("Listar paquetes instalados ->")
	color.Green("upkg installed")
}
