// cmd/check.go
package cmd

import "github.com/spf13/cobra"

// checkCmd representa el comando padre 'check' que agrupará otros subcomandos.
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Ejecuta un conjunto de validaciones específicas",
	Long:  `El comando 'check' agrupa todos los análisis individuales que se pueden ejecutar de forma no interactiva.`,
}

func init() {
	// Añadimos el comando 'check' al comando raíz 'kuma-doctor'
	rootCmd.AddCommand(checkCmd)
}
