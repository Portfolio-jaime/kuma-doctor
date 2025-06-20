// internal/tui/menu.go
package tui

import (
	"fmt"
	"kuma-doctor/internal/kubernetes"
	"kuma-doctor/internal/report"
	"kuma-doctor/pkg/analysis"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

// ShowInteractiveMenu muestra el menú principal y maneja la selección del usuario.
func ShowInteractiveMenu(outputFormat, outputFile string) error {
	for {
		choice := ""
		prompt := &survey.Select{
			Message: "¿Qué aspecto de Kuma Mesh deseas analizar?",
			Options: []string{
				"Resumen General de Salud",
				"Estado de todos los Dataplanes (Proxies)",
				"Consistencia de Políticas de Tráfico (MeshTrafficPermission)",
				"Configuración de mTLS (Seguridad)",
				"Políticas de Resiliencia (Retries, Timeouts, etc.)",
				"Salir",
			},
			PageSize: 10,
		}
		err := survey.AskOne(prompt, &choice)
		if err != nil {
			fmt.Println("Operación cancelada.")
			return nil
		}

		switch choice {
		case "Resumen General de Salud":
			handleSummaryAnalysis(outputFormat, outputFile)
		case "Estado de todos los Dataplanes (Proxies)":
			handleDataplaneAnalysis(outputFormat, outputFile)
		case "Consistencia de Políticas de Tráfico (MeshTrafficPermission)":
			handleTrafficPermissionAnalysis(outputFormat, outputFile)
		case "Configuración de mTLS (Seguridad)":
			handleMTLSAnalysis(outputFormat, outputFile)
		case "Políticas de Resiliencia (Retries, Timeouts, etc.)":
			handleResilienceAnalysis(outputFormat, outputFile)
		case "Salir":
			fmt.Println("¡Hasta luego!")
			return nil
		}
		fmt.Println("\n---\n") // Separador para la siguiente acción
	}
}

// handleSummaryAnalysis ejecuta la lógica para el Resumen General.
func handleSummaryAnalysis(outputFormat, outputFile string) {
	fmt.Println("Generando resumen de salud del mesh...")
	client, err := kubernetes.NewClient()
	if err != nil {
		fmt.Printf("Error al conectar con Kubernetes: %v\n", err)
		return
	}

	result, err := analysis.AnalyzeSummary(client)
	if err != nil {
		fmt.Printf("Error durante el análisis: %v\n", err)
		return
	}

	generateAndDisplayReport(result, outputFormat, outputFile)
}

// handleDataplaneAnalysis ejecuta la lógica para el análisis de Dataplanes.
func handleDataplaneAnalysis(outputFormat, outputFile string) {
	fmt.Println("Analizando Dataplanes...")
	client, err := kubernetes.NewClient()
	if err != nil {
		fmt.Printf("Error al conectar con Kubernetes: %v\n", err)
		return
	}

	result, err := analysis.AnalyzeDataplanes(client)
	if err != nil {
		fmt.Printf("Error durante el análisis: %v\n", err)
		return
	}

	generateAndDisplayReport(result, outputFormat, outputFile)
}

// handleTrafficPermissionAnalysis ejecuta la lógica para el análisis de MTPs.
func handleTrafficPermissionAnalysis(outputFormat, outputFile string) {
	fmt.Println("Analizando consistencia de MeshTrafficPermissions...")
	client, err := kubernetes.NewClient()
	if err != nil {
		fmt.Printf("Error al conectar con Kubernetes: %v\n", err)
		return
	}

	result, err := analysis.AnalyzeTrafficPermissions(client)
	if err != nil {
		fmt.Printf("Error durante el análisis: %v\n", err)
		return
	}

	generateAndDisplayReport(result, outputFormat, outputFile)
}

// handleMTLSAnalysis ejecuta la lógica para el análisis de mTLS.
func handleMTLSAnalysis(outputFormat, outputFile string) {
	fmt.Println("Analizando configuración de mTLS...")
	client, err := kubernetes.NewClient()
	if err != nil {
		fmt.Printf("Error al conectar con Kubernetes: %v\n", err)
		return
	}

	result, err := analysis.AnalyzeMTLS(client)
	if err != nil {
		fmt.Printf("Error durante el análisis: %v\n", err)
		return
	}

	generateAndDisplayReport(result, outputFormat, outputFile)
}

// handleResilienceAnalysis ejecuta la lógica para el análisis de políticas de resiliencia.
func handleResilienceAnalysis(outputFormat, outputFile string) {
	fmt.Println("Analizando políticas de resiliencia...")
	client, err := kubernetes.NewClient()
	if err != nil {
		fmt.Printf("Error al conectar con Kubernetes: %v\n", err)
		return
	}

	result, err := analysis.AnalyzeResilience(client)
	if err != nil {
		fmt.Printf("Error durante el análisis: %v\n", err)
		return
	}

	generateAndDisplayReport(result, outputFormat, outputFile)
}

// generateAndDisplayReport es una función de ayuda para evitar duplicar código.
// Toma el resultado del análisis y los flags de output para generar y mostrar/guardar el reporte.
func generateAndDisplayReport(result *analysis.ValidationResult, outputFormat, outputFile string) {
	if result == nil || len(result.Findings) == 0 {
		fmt.Println("No se generaron hallazgos durante el análisis.")
		return
	}

	reporter, err := report.GetReporter(outputFormat)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	output, err := reporter.Generate(result)
	if err != nil {
		fmt.Printf("Error al generar el reporte: %v\n", err)
		return
	}

	if outputFile != "" {
		err = os.WriteFile(outputFile, []byte(output), 0644)
		if err != nil {
			fmt.Printf("Error al escribir el archivo: %v\n", err)
		} else {
			fmt.Printf("Reporte guardado en %s\n", outputFile)
		}
	} else {
		fmt.Println(output)
	}
}
