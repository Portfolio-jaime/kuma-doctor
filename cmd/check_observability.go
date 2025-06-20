// cmd/check_observability.go
package cmd

import (
	"fmt"
	"kuma-doctor/internal/kubernetes"
	"kuma-doctor/internal/report"
	"kuma-doctor/pkg/analysis"
	"os"

	"github.com/spf13/cobra"
)

var checkObservabilityCmd = &cobra.Command{
	Use:     "observability",
	Short:   "Revisa la cobertura de políticas de observabilidad (Log, Metric, Trace)",
	Aliases: []string{"obs"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Analizando políticas de observabilidad...")
		client, err := kubernetes.NewClient()
		if err != nil {
			fmt.Printf("Error al conectar con Kubernetes: %v\n", err)
			os.Exit(1)
		}

		result, err := analysis.AnalyzeObservability(client)
		if err != nil {
			fmt.Printf("Error durante el análisis: %v\n", err)
			os.Exit(1)
		}

		reporter, err := report.GetReporter(outputFormat)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		output, err := reporter.Generate(result)
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
	checkCmd.AddCommand(checkObservabilityCmd)
}
