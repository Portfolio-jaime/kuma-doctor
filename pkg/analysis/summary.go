// pkg/analysis/summary.go
package analysis

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// AnalyzeSummary ejecuta un análisis de alto nivel de todo el mesh.
func AnalyzeSummary(client dynamic.Interface) (*ValidationResult, error) {
	summary := SummaryStatus{}

	// 1. Contar Meshes
	meshGVR := schema.GroupVersionResource{Group: "kuma.io", Version: "v1alpha1", Resource: "meshes"}
	meshes, err := client.Resource(meshGVR).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error al listar Meshes: %w", err)
	}
	summary.TotalMeshes = len(meshes.Items)

	// 2. Contar y clasificar Dataplanes
	dataplaneGVR := schema.GroupVersionResource{Group: "kuma.io", Version: "v1alpha1", Resource: "dataplanes"}
	dataplanes, err := client.Resource(dataplaneGVR).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error al listar Dataplanes: %w", err)
	}
	summary.TotalDataplanes = len(dataplanes.Items)

	for _, dp := range dataplanes.Items {
		status, _ := getDataplaneStatusFromInbounds(dp)
		switch status.Overall {
		case "Online":
			summary.OnlineDataplanes++
		case "Offline":
			summary.OfflineDataplanes++
		case "Degraded":
			summary.DegradedDataplanes++
		case "Info":
			summary.InfoDataplanes++
		}
	}

	// 3. Contar Políticas (ejemplo con MeshTrafficPermission)
	mtpGVR := schema.GroupVersionResource{Group: "kuma.io", Version: "v1alpha1", Resource: "meshtrafficpermissions"}
	policies, err := client.Resource(mtpGVR).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		// No hacemos que falle todo si solo falla un tipo de política
		fmt.Printf("Advertencia: no se pudieron listar MeshTrafficPermissions: %v\n", err)
	} else {
		summary.TotalPolicies = len(policies.Items)
	}

	result := &ValidationResult{
		Title:       "Resumen General de Salud del Mesh",
		GeneratedAt: time.Now(),
		Findings:    []interface{}{summary}, // El único hallazgo es el propio resumen
	}

	return result, nil
}
