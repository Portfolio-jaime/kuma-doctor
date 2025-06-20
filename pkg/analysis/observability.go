// pkg/analysis/observability.go
package analysis

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// AnalyzeObservability revisa la configuración de políticas como MeshLog, MeshMetric, etc.
func AnalyzeObservability(client dynamic.Interface) (*ValidationResult, error) {
	var findings []interface{}

	// 1. Analizar MeshLog
	logGVR := schema.GroupVersionResource{Group: "kuma.io", Version: "v1alpha1", Resource: "meshlogs"}
	logPolicies, err := client.Resource(logGVR).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error al listar MeshLogs: %w", err)
	}
	if len(logPolicies.Items) == 0 {
		findings = append(findings, ObservabilityFinding{
			Level:      "WARN",
			PolicyType: "MeshLog",
			Resource:   "Global",
			Message:    "No se encontró ninguna política MeshLog. Los logs de acceso no están siendo capturados.",
		})
	} else {
		for _, policy := range logPolicies.Items {
			findings = append(findings, ObservabilityFinding{
				Level:      "INFO",
				PolicyType: "MeshLog",
				Resource:   policy.GetName(),
				Message:    "Política de logging encontrada.",
			})
		}
	}

	// 2. Analizar MeshMetric
	metricGVR := schema.GroupVersionResource{Group: "kuma.io", Version: "v1alpha1", Resource: "meshmetrics"}
	metricPolicies, err := client.Resource(metricGVR).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error al listar MeshMetrics: %w", err)
	}
	if len(metricPolicies.Items) == 0 {
		findings = append(findings, ObservabilityFinding{
			Level:      "WARN",
			PolicyType: "MeshMetric",
			Resource:   "Global",
			Message:    "No se encontró ninguna política MeshMetric. Las métricas para Prometheus pueden no estar habilitadas.",
		})
	} else {
		for _, policy := range metricPolicies.Items {
			findings = append(findings, ObservabilityFinding{
				Level:      "INFO",
				PolicyType: "MeshMetric",
				Resource:   policy.GetName(),
				Message:    "Política de métricas encontrada.",
			})
		}
	}

	// 3. Analizar MeshTrace
	traceGVR := schema.GroupVersionResource{Group: "kuma.io", Version: "v1alpha1", Resource: "meshtraces"}
	tracePolicies, err := client.Resource(traceGVR).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error al listar MeshTraces: %w", err)
	}
	if len(tracePolicies.Items) == 0 {
		findings = append(findings, ObservabilityFinding{
			Level:      "WARN",
			PolicyType: "MeshTrace",
			Resource:   "Global",
			Message:    "No se encontró ninguna política MeshTrace. El tracing distribuido puede no estar configurado.",
		})
	} else {
		for _, policy := range tracePolicies.Items {
			findings = append(findings, ObservabilityFinding{
				Level:      "INFO",
				PolicyType: "MeshTrace",
				Resource:   policy.GetName(),
				Message:    "Política de tracing encontrada.",
			})
		}
	}

	return &ValidationResult{
		Title:       "Análisis de Políticas de Observabilidad",
		GeneratedAt: time.Now(),
		Findings:    findings,
	}, nil
}
