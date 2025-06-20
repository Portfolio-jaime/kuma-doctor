// pkg/analysis/dataplanes.go
package analysis

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// AnalyzeDataplanes ejecuta la validación de todos los dataplanes y devuelve un resultado estructurado.
func AnalyzeDataplanes(client dynamic.Interface) (*ValidationResult, error) {
	dataplaneGVR := schema.GroupVersionResource{
		Group:    "kuma.io",
		Version:  "v1alpha1",
		Resource: "dataplanes",
	}

	unstructuredDataplanes, err := client.Resource(dataplaneGVR).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error al listar Dataplanes: %w", err)
	}

	result := &ValidationResult{
		Title:       "Análisis de Estado de Dataplanes",
		GeneratedAt: time.Now(),
	}

	for _, dp := range unstructuredDataplanes.Items {
		status, _ := getDataplaneStatusFromInbounds(dp)
		finding := DataplaneStatus{
			Name:      dp.GetName(),
			Namespace: dp.GetNamespace(),
			Status:    status.Overall,
			Details:   status.Details,
		}
		result.Findings = append(result.Findings, finding)
	}

	return result, nil
}

// Estructura interna para el estado derivado de los inbounds
type derivedStatus struct {
	Overall string
	Details string
}

func getDataplaneStatusFromInbounds(dp unstructured.Unstructured) (derivedStatus, []string) {
	inbounds, found, err := unstructured.NestedSlice(dp.Object, "spec", "networking", "inbound")
	if err != nil || !found || len(inbounds) == 0 {
		return derivedStatus{"Info", "Sin inbounds definidos"}, nil
	}

	totalInbounds := len(inbounds)
	readyInbounds := 0
	var unhealthyDetails []string

	for _, inboundItem := range inbounds {
		inboundMap, ok := inboundItem.(map[string]interface{})
		if !ok {
			continue
		}
		isReady, readyFound, _ := unstructured.NestedBool(inboundMap, "health", "ready")
		if readyFound && isReady {
			readyInbounds++
		} else {
			port, _, _ := unstructured.NestedInt64(inboundMap, "port")
			service, _, _ := unstructured.NestedString(inboundMap, "tags", "kuma.io/service")
			unhealthyDetails = append(unhealthyDetails, fmt.Sprintf("Puerto %d (Servicio: %s) no está 'ready'", port, service))
		}
	}

	if readyInbounds == totalInbounds {
		return derivedStatus{"Online", "Todos los inbounds 'ready'"}, nil
	} else if readyInbounds > 0 {
		return derivedStatus{"Degraded", fmt.Sprintf("%d de %d inbounds 'ready'", readyInbounds, totalInbounds)}, unhealthyDetails
	} else {
		return derivedStatus{"Offline", "Ningún inbound está 'ready'"}, unhealthyDetails
	}
}
