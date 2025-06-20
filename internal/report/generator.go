// internal/report/generator.go
package report

import (
	"encoding/json"
	"fmt"
	"kuma-doctor/pkg/analysis"
	"strings"
	"text/tabwriter"
	"time"
)

// Reporter define la interfaz para todos los generadores de reportes.
type Reporter interface {
	Generate(result *analysis.ValidationResult) (string, error)
}

// GetReporter es una factory para obtener el reporter correcto según el formato.
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

func (r *TextReporter) Generate(result *analysis.ValidationResult) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("--- %s ---\n", result.Title))
	sb.WriteString(fmt.Sprintf("Fecha: %s\n\n", result.GeneratedAt.Format(time.RFC1123)))

	if len(result.Findings) == 0 {
		sb.WriteString("No se encontraron hallazgos.\n")
		return sb.String(), nil
	}

	w := tabwriter.NewWriter(&sb, 0, 0, 3, ' ', 0)

	// Verificamos el tipo de hallazgo para formatearlo correctamente
	switch result.Findings[0].(type) {
	case analysis.DataplaneStatus:
		fmt.Fprintln(w, "NOMBRE\tNAMESPACE\tESTADO\tDETALLES")
		fmt.Fprintln(w, "------\t---------\t------\t--------")
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
			fmt.Fprintf(w, "%s\t%s\t%s %s\t%s\n", dpStatus.Name, dpStatus.Namespace, emoji, dpStatus.Status, dpStatus.Details)
		}

	case analysis.SummaryStatus:
		summary := result.Findings[0].(analysis.SummaryStatus)
		fmt.Fprintln(w, "RECURSO\tCANTIDAD\t")
		fmt.Fprintln(w, "-------\t--------\t")
		fmt.Fprintf(w, "Meshes\t%d\t\n", summary.TotalMeshes)
		fmt.Fprintln(w, "\t\t") // Separador visual
		fmt.Fprintf(w, "Dataplanes Totales\t%d\t\n", summary.TotalDataplanes)
		fmt.Fprintf(w, "  ✅ En Línea\t%d\t\n", summary.OnlineDataplanes)
		fmt.Fprintf(w, "  ❌ Fuera de Línea\t%d\t\n", summary.OfflineDataplanes)
		fmt.Fprintf(w, "  ⚠️ Degradados\t%d\t\n", summary.DegradedDataplanes)
		fmt.Fprintf(w, "  ℹ️ Informativos\t%d\t\n", summary.InfoDataplanes)
		fmt.Fprintln(w, "\t\t")
		fmt.Fprintf(w, "Políticas de Tráfico (MTPs)\t%d\t\n", summary.TotalPolicies)

	case analysis.PolicyFinding:
		fmt.Fprintln(w, "NIVEL\tRECURSO\tMENSAJE")
		fmt.Fprintln(w, "-----\t-------\t-------")
		for _, finding := range result.Findings {
			policyFinding, _ := finding.(analysis.PolicyFinding)
			var emoji string
			switch policyFinding.Level {
			case "ALERT":
				emoji = "🚨"
			case "WARN":
				emoji = "⚠️"
			case "INFO":
				emoji = "ℹ️"
			}
			message := strings.ReplaceAll(policyFinding.Message, "\t", " ")
			fmt.Fprintf(w, "%s %s\t%s\t%s\n", emoji, policyFinding.Level, policyFinding.Resource, message)
		}

	case analysis.MTLSFinding:
		fmt.Fprintln(w, "NIVEL\tRECURSO\tMENSAJE")
		fmt.Fprintln(w, "-----\t-------\t-------")
		for _, finding := range result.Findings {
			mtlsFinding, _ := finding.(analysis.MTLSFinding)
			var emoji string
			switch mtlsFinding.Level {
			case "ALERT":
				emoji = "🚨"
			case "WARN":
				emoji = "⚠️"
			case "INFO":
				emoji = "✅"
			}
			fmt.Fprintf(w, "%s %s\t%s\t%s\n", emoji, mtlsFinding.Level, mtlsFinding.Resource, mtlsFinding.Message)
		}

	case analysis.ResilienceFinding:
		fmt.Fprintln(w, "NIVEL\tTIPO DE POLÍTICA\tSERVICIO\tMENSAJE")
		fmt.Fprintln(w, "-----\t----------------\t--------\t-------")
		for _, finding := range result.Findings {
			resilienceFinding, _ := finding.(analysis.ResilienceFinding)
			var emoji string
			switch resilienceFinding.Level {
			case "WARN":
				emoji = "⚠️"
			case "INFO":
				emoji = "✅"
			}
			fmt.Fprintf(w, "%s %s\t%s\t%s\t%s\n", emoji, resilienceFinding.Level, resilienceFinding.PolicyType, resilienceFinding.Service, resilienceFinding.Message)
		}
	}

	w.Flush()
	return sb.String(), nil
}

// --- Implementación de JsonReporter ---
type JsonReporter struct{}

func (r *JsonReporter) Generate(result *analysis.ValidationResult) (string, error) {
	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// --- Implementación de MarkdownReporter ---
type MarkdownReporter struct{}

func (r *MarkdownReporter) Generate(result *analysis.ValidationResult) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n\n", result.Title))
	sb.WriteString(fmt.Sprintf("**Fecha:** %s\n\n", result.GeneratedAt.Format(time.RFC1123)))

	if len(result.Findings) == 0 {
		sb.WriteString("No se encontraron hallazgos.\n")
		return sb.String(), nil
	}

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
		sb.WriteString("## Resumen del Mesh\n\n")
		sb.WriteString(fmt.Sprintf("- **Meshes:** %d\n", summary.TotalMeshes))
		sb.WriteString(fmt.Sprintf("- **Dataplanes Totales:** %d\n", summary.TotalDataplanes))
		sb.WriteString(fmt.Sprintf("  - ✅ **En Línea:** %d\n", summary.OnlineDataplanes))
		sb.WriteString(fmt.Sprintf("  - ❌ **Fuera de Línea:** %d\n", summary.OfflineDataplanes))
		sb.WriteString(fmt.Sprintf("  - ⚠️ **Degradados:** %d\n", summary.DegradedDataplanes))
		sb.WriteString(fmt.Sprintf("  - ℹ️ **Informativos:** %d\n", summary.InfoDataplanes))
		sb.WriteString(fmt.Sprintf("- **Políticas de Tráfico (MTPs):** %d\n", summary.TotalPolicies))

	case analysis.PolicyFinding:
		sb.WriteString("| Nivel | Recurso | Mensaje |\n")
		sb.WriteString("|---|---|---|\n")
		for _, finding := range result.Findings {
			policyFinding, _ := finding.(analysis.PolicyFinding)
			var emoji string
			switch policyFinding.Level {
			case "ALERT":
				emoji = "🚨"
			case "WARN":
				emoji = "⚠️"
			case "INFO":
				emoji = "ℹ️"
			}
			sb.WriteString(fmt.Sprintf("| %s %s | `%s` | %s |\n", emoji, policyFinding.Level, policyFinding.Resource, policyFinding.Message))
		}

	case analysis.MTLSFinding:
		sb.WriteString("| Nivel | Recurso | Mensaje |\n")
		sb.WriteString("|---|---|---|\n")
		for _, finding := range result.Findings {
			mtlsFinding, _ := finding.(analysis.MTLSFinding)
			var emoji string
			switch mtlsFinding.Level {
			case "ALERT":
				emoji = "🚨"
			case "WARN":
				emoji = "⚠️"
			case "INFO":
				emoji = "✅"
			}
			sb.WriteString(fmt.Sprintf("| %s %s | `%s` | %s |\n", emoji, mtlsFinding.Level, mtlsFinding.Resource, mtlsFinding.Message))
		}

	case analysis.ResilienceFinding:
		sb.WriteString("| Nivel | Tipo de Política | Servicio | Mensaje |\n")
		sb.WriteString("|---|---|---|---|\n")
		for _, finding := range result.Findings {
			resilienceFinding, _ := finding.(analysis.ResilienceFinding)
			var emoji string
			switch resilienceFinding.Level {
			case "WARN":
				emoji = "⚠️"
			case "INFO":
				emoji = "✅"
			}
			sb.WriteString(fmt.Sprintf("| %s %s | `%s` | `%s` | %s |\n", emoji, resilienceFinding.Level, resilienceFinding.PolicyType, resilienceFinding.Service, resilienceFinding.Message))
		}
	}

	return sb.String(), nil
}
