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
	// Add your code to evaluate all existing results here
	// Specify the directory containing JSON files
	directory := "../test-web-page/results/"

	// List files in the directory
	files, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Iterate over files
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

	// Parse the input date string to verify if it's in correct format
	_, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		fmt.Println("Formato de fecha incorrecto. Use el formato YYYY-MM-DD.")
		return
	}

	fmt.Printf("Evaluando por fecha %s...\n", fecha)

	// Specify the directory containing JSON files
	directory := "../test-web-page/results/"

	// List files in the directory
	files, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Iterate over files
	for _, file := range files {
		// Check if the file has .json extension
		if strings.HasSuffix(file.Name(), ".json") {
			// Extract date from filename
			fileDate := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			if strings.Contains(fileDate, fecha) {
				// If the extracted date matches the provided date, read and process the file
				ReadJson(filepath.Join(directory, file.Name()))
			}
		}
	}
}
