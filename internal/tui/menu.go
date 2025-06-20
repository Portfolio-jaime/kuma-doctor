// internal/tui/menu.go
package tui

import (
	"fmt"
	"kuma-doctor/internal/kubernetes"
	"kuma-doctor/internal/report"
	"kuma-doctor/pkg/analysis"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"k8s.io/client-go/dynamic"
)

// ShowInteractiveMenu muestra el menú principal y maneja la selección del usuario.
func ShowInteractiveMenu(outputFormat, outputFile string) error {
	for {
		choice := ""
		prompt := &survey.Select{
			Message: "¿Qué aspecto de Kuma Mesh deseas analizar?",
			Options: []string{
				"Generar Reporte Completo", // <-- Nueva opción funcional
				"Resumen General de Salud",
				"Estado de todos los Dataplanes (Proxies)",
				"Consistencia de Políticas de Tráfico (MeshTrafficPermission)",
				"Configuración de mTLS (Seguridad)",
				"Políticas de Resiliencia (Retries, Timeouts, etc.)",
				"Políticas de Observabilidad (Logs, Metrics, Traces)",
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
		case "Generar Reporte Completo":
			handleFullReportAnalysis(outputFormat, outputFile)
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
		case "Políticas de Observabilidad (Logs, Metrics, Traces)":
			handleObservabilityAnalysis(outputFormat, outputFile)
		case "Salir":
			fmt.Println("¡Hasta luego!")
			return nil
		}
		fmt.Println("\n---\n")
	}
}

// --- Nueva función para el Reporte Completo ---
func handleFullReportAnalysis(outputFormat, outputFile string) {
	fmt.Println("Generando reporte completo, esto puede tardar un momento...")
	client, err := kubernetes.NewClient()
	if err != nil {
		fmt.Printf("Error al conectar con Kubernetes: %v\n", err)
		return
	}

	// Creamos una lista para almacenar todos los resultados
	var allResults []*analysis.ValidationResult

	// Ejecutamos cada análisis y añadimos su resultado a la lista
	if result, err := analysis.AnalyzeSummary(client); err == nil {
		allResults = append(allResults, result)
	}
	if result, err := analysis.AnalyzeDataplanes(client); err == nil {
		allResults = append(allResults, result)
	}
	if result, err := analysis.AnalyzeTrafficPermissions(client); err == nil {
		allResults = append(allResults, result)
	}
	if result, err := analysis.AnalyzeMTLS(client); err == nil {
		allResults = append(allResults, result)
	}
	if result, err := analysis.AnalyzeResilience(client); err == nil {
		allResults = append(allResults, result)
	}
	if result, err := analysis.AnalyzeObservability(client); err == nil {
		allResults = append(allResults, result)
	}

	// Pasamos la lista completa al generador de reportes
	generateAndDisplayReport(allResults, outputFormat, outputFile)
}

// --- Handlers Anteriores (Ahora envuelven el resultado en una lista) ---

func handleSummaryAnalysis(outputFormat, outputFile string) {
	fmt.Println("Generando resumen de salud del mesh...")
	executeAnalysis(analysis.AnalyzeSummary, outputFormat, outputFile)
}
func handleDataplaneAnalysis(outputFormat, outputFile string) {
	fmt.Println("Analizando Dataplanes...")
	executeAnalysis(analysis.AnalyzeDataplanes, outputFormat, outputFile)
}
func handleTrafficPermissionAnalysis(outputFormat, outputFile string) {
	fmt.Println("Analizando consistencia de MeshTrafficPermissions...")
	executeAnalysis(analysis.AnalyzeTrafficPermissions, outputFormat, outputFile)
}
func handleMTLSAnalysis(outputFormat, outputFile string) {
	fmt.Println("Analizando configuración de mTLS...")
	executeAnalysis(analysis.AnalyzeMTLS, outputFormat, outputFile)
}
func handleResilienceAnalysis(outputFormat, outputFile string) {
	fmt.Println("Analizando políticas de resiliencia...")
	executeAnalysis(analysis.AnalyzeResilience, outputFormat, outputFile)
}
func handleObservabilityAnalysis(outputFormat, outputFile string) {
	fmt.Println("Analizando políticas de observabilidad...")
	executeAnalysis(analysis.AnalyzeObservability, outputFormat, outputFile)
}

// --- Funciones Helper (Actualizadas para el nuevo Reporter) ---

func executeAnalysis(analysisFunc func(client dynamic.Interface) (*analysis.ValidationResult, error), outputFormat, outputFile string) {
	client, err := kubernetes.NewClient()
	if err != nil {
		fmt.Printf("Error al conectar con Kubernetes: %v\n", err)
		return
	}
	result, err := analysisFunc(client)
	if err != nil {
		fmt.Printf("Error durante el análisis: %v\n", err)
		return
	}
	// Envolvemos el resultado único en una lista para usar la nueva interfaz del reporter
	generateAndDisplayReport([]*analysis.ValidationResult{result}, outputFormat, outputFile)
}

func generateAndDisplayReport(results []*analysis.ValidationResult, outputFormat, outputFile string) {
	if len(results) == 0 {
		color.Green("✅ No se generaron hallazgos durante el análisis.")
		return
	}

	reporter, err := report.GetReporter(outputFormat)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	output, err := reporter.Generate(results)
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
