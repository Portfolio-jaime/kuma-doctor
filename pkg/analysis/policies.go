// pkg/analysis/policies.go
package analysis

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// AnalyzeTrafficPermissions revisa la configuración y consistencia de MeshTrafficPermissions.
func AnalyzeTrafficPermissions(client dynamic.Interface) (*ValidationResult, error) {
	// GVRs para los recursos que necesitamos
	mtpGVR := schema.GroupVersionResource{Group: "kuma.io", Version: "v1alpha1", Resource: "meshtrafficpermissions"}
	dataplaneGVR := schema.GroupVersionResource{Group: "kuma.io", Version: "v1alpha1", Resource: "dataplanes"}

	// 1. Obtener todas las políticas y todos los dataplanes
	policies, err := client.Resource(mtpGVR).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error al listar MeshTrafficPermissions: %w", err)
	}
	dataplanes, err := client.Resource(dataplaneGVR).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error al listar Dataplanes: %w", err)
	}

	// 2. Construir un mapa de todos los servicios existentes y los servicios protegidos por políticas
	allServices := make(map[string]bool)
	protectedServices := make(map[string]bool)
	var findings []interface{}

	for _, dp := range dataplanes.Items {
		// El nombre del servicio se encuentra en las etiquetas del inbound
		inbounds, found, _ := unstructured.NestedSlice(dp.Object, "spec", "networking", "inbound")
		if found {
			for _, inboundItem := range inbounds {
				inboundMap, _ := inboundItem.(map[string]interface{})
				serviceName, _, _ := unstructured.NestedString(inboundMap, "tags", "kuma.io/service")
				if serviceName != "" {
					allServices[serviceName] = true
				}
			}
		}
	}

	// 3. Analizar cada política
	for _, policy := range policies.Items {
		// Analizar a quién protege la política (sección `to`)
		to, found, _ := unstructured.NestedSlice(policy.Object, "spec", "to")
		if found && len(to) > 0 {
			target, _ := to[0].(map[string]interface{}) // Asumimos una sola regla 'to' por simplicidad
			targetRef, _, _ := unstructured.NestedString(target, "targetRef", "name")
			if targetRef == "*" { // Si la política apunta a todos los servicios
				for service := range allServices {
					protectedServices[service] = true
				}
			} else {
				protectedServices[targetRef] = true
			}
		}

		// Analizar quién tiene permiso (sección `from`) para alertas de seguridad
		from, found, _ := unstructured.NestedSlice(policy.Object, "spec", "from")
		if found && len(from) > 0 {
			source, _ := from[0].(map[string]interface{})
			sourceRef, _, _ := unstructured.NestedString(source, "targetRef", "name")
			if sourceRef == "*" {
				findings = append(findings, PolicyFinding{
					Level:    "INFO",
					Message:  "La política permite tráfico desde CUALQUIER servicio. Asegúrate de que esto sea intencional.",
					Resource: policy.GetName(),
				})
			}
		}
	}

	// 4. Comparar todos los servicios con los servicios protegidos
	for service := range allServices {
		if !protectedServices[service] {
			findings = append(findings, PolicyFinding{
				Level:    "ALERT",
				Message:  "Este servicio no está protegido por ninguna MeshTrafficPermission. Podría estar aislado si la política por defecto es 'deny'.",
				Resource: service,
			})
		}
	}

	if len(findings) == 0 {
		findings = append(findings, PolicyFinding{
			Level:    "INFO",
			Message:  "Todos los servicios están cubiertos por al menos una MeshTrafficPermission.",
			Resource: "Global",
		})
	}

	return &ValidationResult{
		Title:       "Análisis de Consistencia de Políticas de Tráfico",
		GeneratedAt: time.Now(),
		Findings:    findings,
	}, nil
}
