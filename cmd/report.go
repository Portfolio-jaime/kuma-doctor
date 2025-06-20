// cmd/report.go
package cmd

import (
	"fmt"
	"kuma-doctor/internal/kubernetes"
	"kuma-doctor/internal/report"
	"kuma-doctor/pkg/analysis"
	"os"

	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Genera un reporte completo con todos los análisis disponibles",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generando reporte completo, esto puede tardar un momento...")
		client, err := kubernetes.NewClient()
		if err != nil {
			fmt.Printf("Error al conectar con Kubernetes: %v\n", err)
			os.Exit(1)
		}

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

		reporter, err := report.GetReporter(outputFormat)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		output, err := reporter.Generate(allResults)
		if err != nil {
			fmt.Printf("Error al generar el reporte: %v\n", err)
			os.Exit(1)
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
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}
