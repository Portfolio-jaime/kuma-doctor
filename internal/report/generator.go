// internal/report/generator.go
package report

import (
	"encoding/json"
	"fmt"
	"kuma-doctor/pkg/analysis"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
)

// --- Definición de Colores para la Consola ---
var (
	green  = color.New(color.FgGreen).SprintfFunc()
	red    = color.New(color.FgRed).SprintfFunc()
	yellow = color.New(color.FgYellow).SprintfFunc()
	cyan   = color.New(color.FgCyan).SprintfFunc()
	bold   = color.New(color.Bold).SprintFunc()
)

// Reporter define la interfaz para todos los generadores de reportes.
type Reporter interface {
	Generate(results []*analysis.ValidationResult) (string, error)
}

// GetReporter es una factory para obtener el reporter correcto según el formato.
// ESTA ES LA FUNCIÓN CORREGIDA Y COMPLETA.
func GetReporter(format string) (Reporter, error) {
	switch format {
	case "txt":
		return &TextReporter{}, nil
	case "md":
		return &MarkdownReporter{}, nil
	case "json":
		return &JsonReporter{}, nil
	default:
		return nil, fmt.Errorf("formato de reporte desconocido: %s", format)
	}
}

// --- Implementación de TextReporter ---
type TextReporter struct{}

func (r *TextReporter) Generate(results []*analysis.ValidationResult) (string, error) {
	var finalReport strings.Builder
	for i, result := range results {
		var sb strings.Builder
		sb.WriteString(bold(fmt.Sprintf("--- %s ---\n", result.Title)))
		sb.WriteString(fmt.Sprintf("Fecha: %s\n\n", result.GeneratedAt.Format(time.RFC1123)))

		if len(result.Findings) == 0 {
			sb.WriteString(green("✅ No se encontraron hallazgos problemáticos.\n"))
		} else {
			w := tabwriter.NewWriter(&sb, 0, 0, 3, ' ', 0)
			switch result.Findings[0].(type) {
			case analysis.DataplaneStatus:
				fmt.Fprintln(w, bold("NOMBRE\tNAMESPACE\tESTADO\tDETALLES"))
				fmt.Fprintln(w, bold("------\t---------\t------\t--------"))
				for _, finding := range result.Findings {
					dpStatus, _ := finding.(analysis.DataplaneStatus)
					var statusCell string
					switch dpStatus.Status {
					case "Online":
						statusCell = green("✅ Online")
					case "Offline":
						statusCell = red("❌ Offline")
					case "Degraded":
						statusCell = yellow("⚠️ Degraded")
					case "Info":
						statusCell = cyan("ℹ️ Info")
					}
					fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", dpStatus.Name, dpStatus.Namespace, statusCell, dpStatus.Details)
				}
			case analysis.SummaryStatus:
				summary := result.Findings[0].(analysis.SummaryStatus)
				fmt.Fprintln(w, bold("RECURSO\tCANTIDAD\t"))
				fmt.Fprintln(w, bold("-------\t--------\t"))
				fmt.Fprintf(w, "%s\t%d\t\n", "Meshes", summary.TotalMeshes)
				fmt.Fprintln(w, "\t\t")
				fmt.Fprintf(w, "%s\t%d\t\n", "Dataplanes Totales", summary.TotalDataplanes)
				fmt.Fprintf(w, "  %s\t%d\t\n", green("✅ En Línea"), summary.OnlineDataplanes)
				fmt.Fprintf(w, "  %s\t%d\t\n", red("❌ Fuera de Línea"), summary.OfflineDataplanes)
				fmt.Fprintf(w, "  %s\t%d\t\n", yellow("⚠️ Degradados"), summary.DegradedDataplanes)
				fmt.Fprintf(w, "  %s\t%d\t\n", cyan("ℹ️ Informativos"), summary.InfoDataplanes)
				fmt.Fprintln(w, "\t\t")
				fmt.Fprintf(w, "%s\t%d\t\n", "Políticas de Tráfico (MTPs)", summary.TotalPolicies)
			case analysis.PolicyFinding, analysis.MTLSFinding, analysis.ResilienceFinding, analysis.ObservabilityFinding:
				fmt.Fprintln(w, bold("NIVEL\tRECURSO/TIPO\tMENSAJE"))
				fmt.Fprintln(w, bold("-----\t------------\t-------"))
				for _, finding := range result.Findings {
					var level, resource, message string
					switch f := finding.(type) {
					case analysis.PolicyFinding:
						level, resource, message = f.Level, f.Resource, f.Message
					case analysis.MTLSFinding:
						level, resource, message = f.Level, f.Resource, f.Message
					case analysis.ResilienceFinding:
						level, resource, message = f.Level, f.Service, fmt.Sprintf("(%s) %s", f.PolicyType, f.Message)
					case analysis.ObservabilityFinding:
						level, resource, message = f.Level, f.Resource, fmt.Sprintf("(%s) %s", f.PolicyType, f.Message)
					}
					var levelCell string
					switch level {
					case "ALERT":
						levelCell = red("🚨 ALERT")
					case "WARN":
						levelCell = yellow("⚠️ WARN")
					case "INFO":
						levelCell = green("✅ INFO")
					}
					fmt.Fprintf(w, "%s\t%s\t%s\n", levelCell, resource, message)
				}
			}
			w.Flush()
		}
		finalReport.WriteString(sb.String())
		if i < len(results)-1 {
			finalReport.WriteString("\n\n") // Añade un separador entre reportes
		}
	}
	return finalReport.String(), nil
}

// --- Implementación de JsonReporter ---
type JsonReporter struct{}

func (r *JsonReporter) Generate(results []*analysis.ValidationResult) (string, error) {
	if len(results) == 1 {
		bytes, err := json.MarshalIndent(results[0], "", "  ")
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}
	bytes, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// --- Implementación de MarkdownReporter ---
type MarkdownReporter struct{}

func (r *MarkdownReporter) Generate(results []*analysis.ValidationResult) (string, error) {
	var finalReport strings.Builder
	for i, result := range results {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("## %s\n\n", result.Title))
		sb.WriteString(fmt.Sprintf("**Fecha:** %s\n\n", result.GeneratedAt.Format(time.RFC1123)))
		if len(result.Findings) == 0 {
			sb.WriteString("✅ No se encontraron hallazgos problemáticos.\n")
		} else {
			switch result.Findings[0].(type) {
			case analysis.DataplaneStatus:
				sb.WriteString("| Nombre | Namespace | Estado | Detalles |\n")
				sb.WriteString("|---|---|---|---|\n")
				for _, finding := range result.Findings {
					dpStatus, _ := finding.(analysis.DataplaneStatus)
					var emoji string
					switch dpStatus.Status {
					case "Online":
						emoji = "✅"
					case "Offline":
						emoji = "❌"
					case "Degraded":
						emoji = "⚠️"
					case "Info":
						emoji = "ℹ️"
					}
					sb.WriteString(fmt.Sprintf("| %s | %s | %s %s | %s |\n", dpStatus.Name, dpStatus.Namespace, emoji, dpStatus.Status, dpStatus.Details))
				}
			case analysis.SummaryStatus:
				summary := result.Findings[0].(analysis.SummaryStatus)
				sb.WriteString(fmt.Sprintf("- **Meshes:** %d\n", summary.TotalMeshes))
				sb.WriteString(fmt.Sprintf("- **Dataplanes Totales:** %d\n", summary.TotalDataplanes))
				sb.WriteString(fmt.Sprintf("  - ✅ **En Línea:** %d\n", summary.OnlineDataplanes))
				sb.WriteString(fmt.Sprintf("  - ❌ **Fuera de Línea:** %d\n", summary.OfflineDataplanes))
				sb.WriteString(fmt.Sprintf("  - ⚠️ **Degradados:** %d\n", summary.DegradedDataplanes))
				sb.WriteString(fmt.Sprintf("  - ℹ️ **Informativos:** %d\n", summary.InfoDataplanes))
				sb.WriteString(fmt.Sprintf("- **Políticas de Tráfico (MTPs):** %d\n", summary.TotalPolicies))

			case analysis.PolicyFinding, analysis.MTLSFinding, analysis.ResilienceFinding, analysis.ObservabilityFinding:
				sb.WriteString("| Nivel | Recurso/Tipo | Mensaje |\n")
				sb.WriteString("|---|---|---|\n")
				for _, finding := range result.Findings {
					var level, resource, message, emoji string
					switch f := finding.(type) {
					case analysis.PolicyFinding:
						level, resource, message = f.Level, f.Resource, f.Message
					case analysis.MTLSFinding:
						level, resource, message = f.Level, f.Resource, f.Message
					case analysis.ResilienceFinding:
						level, resource, message = f.Level, f.Service, fmt.Sprintf("_(%s)_ %s", f.PolicyType, f.Message)
					case analysis.ObservabilityFinding:
						level, resource, message = f.Level, f.Resource, fmt.Sprintf("_(%s)_ %s", f.PolicyType, f.Message)
					}
					switch level {
					case "ALERT":
						emoji = "🚨"
					case "WARN":
						emoji = "⚠️"
					case "INFO":
						emoji = "✅"
					}
					sb.WriteString(fmt.Sprintf("| %s %s | `%s` | %s |\n", emoji, level, resource, message))
				}
			}
		}
		finalReport.WriteString(sb.String())
		if i < len(results)-1 {
			finalReport.WriteString("\n---\n\n") // Separador de Markdown
		}
	}
	return finalReport.String(), nil
}
