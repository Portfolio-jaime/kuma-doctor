// cmd/root.go
package cmd

import (
	"fmt"
	"kuma-doctor/internal/tui"
	"os"

	"github.com/spf13/cobra"
)

var (
	outputFormat string
	outputFile   string
)

var rootCmd = &cobra.Command{
	Use:   "kuma-doctor",
	Short: "kuma-doctor es una herramienta CLI para diagnosticar y validar un Kuma mesh.",
	Long: `Una completa herramienta de diagnóstico que te permite revisar la salud
y la configuración de tu Kuma service mesh de manera interactiva o a través
de subcomandos para la automatización.`,
	// Si se ejecuta 'kuma-doctor' sin subcomandos, mostramos el menú.
	Run: func(cmd *cobra.Command, args []string) {
		// Ignora el error aquí, ya que el menú maneja su propio flujo
		_ = tui.ShowInteractiveMenu(outputFormat, outputFile)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Flags globales para todos los comandos
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "txt", "Formato del reporte (txt, md, json)")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "file", "f", "", "Ruta del archivo para guardar el reporte (opcional)")
}
