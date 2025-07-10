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
	if len(os.Args) < 2 {
		color.Red("¿Tenés la cabeza plana?\n")
		printHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "install":
		if len(os.Args) < 3 {
			color.Red("Uso: upkg install <paquete>\n")
			os.Exit(1)
		}
		installPackage(os.Args[2])
	case "list":
		listPackages()
	case "installed":
		listInstalled()
	case "remove":
		if len(os.Args) < 3 {
			color.Red("Uso: upkg remove <paquete>\n")
			os.Exit(1)
		}
		removePackage(os.Args[2])
	case "help":
		printHelp()
	default:
		color.Red("¿Qué intentaste hacer?: %s\n", os.Args[1])
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

func removePackage(pkgName string) {
	pak, err := local.GetInstalledPkg(pkgName)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if err := pak.Uninstall(); err != nil {
		log.Fatalf("Error desinstalando el paquete: %v", err)
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
	fmt.Print("Instalar paquete desde Mullin -> ")
	color.Green("upkg install <paquete>\n")

	fmt.Print("Listar paquetes disponibles en el Mullin -> ")
	color.Green("upkg list\n")

	fmt.Print("Listar paquetes instalados -> ")
	color.Green("upkg installed\n")

	fmt.Print("Desinstalar un paquete -> ")
	color.Green("upkg remove <paquete>\n")
}
