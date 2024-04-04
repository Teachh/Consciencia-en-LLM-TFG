package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	fmt.Println("Instrucciones:")
	fmt.Println("1- Evaluar todos los resultados existentes")
	fmt.Println("2- Evaluar por fecha (YYYY-MM-DD)")

	var choice int
	fmt.Print("Elige una opción (1-2): ")
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		evaluarTodosResultados()
	case 2:
		evaluarPorFecha()
	default:
		fmt.Println("Opción no válida")
	}
}

func evaluarTodosResultados() {
	directory := "../test-web-page/results/"

	files, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			ReadJson(directory + file.Name())
		}
	}
}

func evaluarPorFecha() {
	var fecha string
	fmt.Print("Ingrese la fecha (YYYY-MM-DD): ")
	fmt.Scanln(&fecha)

	_, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		fmt.Println("Formato de fecha incorrecto. Use el formato YYYY-MM-DD.")
		return
	}

	fmt.Printf("Evaluando por fecha %s...\n", fecha)

	directory := "../test-web-page/results/"

	files, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			fileDate := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			if strings.Contains(fileDate, fecha) {
				ReadJson(filepath.Join(directory, file.Name()))
			}
		}
	}
}
