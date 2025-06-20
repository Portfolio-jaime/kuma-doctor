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
				"Resumen General de Salud",
				"Estado de todos los Dataplanes (Proxies)",
				"Consistencia de Políticas de Tráfico (MeshTrafficPermission)",
				"Configuración de mTLS (Seguridad)",
				"Políticas de Resiliencia (Retries, Timeouts, etc.)",
				"Políticas de Observabilidad (Logs, Metrics, Traces)",
				"Generar Reporte Completo (Próximamente)",
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
		case "Políticas de Observabilidad (Logs, Metrics, Traces)":
			handleObservabilityAnalysis(outputFormat, outputFile)
		case "Generar Reporte Completo (Próximamente)":
			fmt.Println("Esta funcionalidad aún no está implementada.")
		case "Salir":
			fmt.Println("¡Hasta luego!")
			return nil
		}
		fmt.Println("\n---\n")
	}
}

// --- Funciones Handler (Ahora más limpias gracias a la refactorización) ---

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

// --- Funciones Helper ---

// executeAnalysis es una función genérica para ejecutar cualquier tipo de análisis.
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

	generateAndDisplayReport(result, outputFormat, outputFile)
}

// generateAndDisplayReport es una función de ayuda para evitar duplicar código.
func generateAndDisplayReport(result *analysis.ValidationResult, outputFormat, outputFile string) {
	if result == nil || len(result.Findings) == 0 {
		color.Green("✅ No se encontraron hallazgos problemáticos.")
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
