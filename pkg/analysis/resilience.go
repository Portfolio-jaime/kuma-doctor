// pkg/analysis/resilience.go
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

// AnalyzeResilience revisa la cobertura de políticas como MeshRetry, MeshTimeout, etc.
func AnalyzeResilience(client dynamic.Interface) (*ValidationResult, error) {
	var findings []interface{}

	// 1. Obtener todos los servicios únicos desde los Dataplanes
	allServices, err := getAllServices(client)
	if err != nil {
		return nil, err
	}

	// 2. Analizar la cobertura para cada tipo de política de resiliencia
	retryCoveredServices, err := getCoveredServices(client, "meshreries", "MeshRetry")
	if err != nil {
		fmt.Printf("Advertencia: no se pudo analizar MeshRetry: %v\n", err)
	}
	timeoutCoveredServices, err := getCoveredServices(client, "meshtimeouts", "MeshTimeout")
	if err != nil {
		fmt.Printf("Advertencia: no se pudo analizar MeshTimeout: %v\n", err)
	}
	breakerCoveredServices, err := getCoveredServices(client, "meshcircuitbreakers", "MeshCircuitBreaker")
	if err != nil {
		fmt.Printf("Advertencia: no se pudo analizar MeshCircuitBreaker: %v\n", err)
	}

	// 3. Comparar y generar hallazgos
	for service := range allServices {
		if _, found := retryCoveredServices[service]; !found {
			findings = append(findings, ResilienceFinding{
				Level:      "WARN",
				PolicyType: "MeshRetry",
				Service:    service,
				Message:    "El servicio no está cubierto por ninguna política de reintentos.",
			})
		}
		if _, found := timeoutCoveredServices[service]; !found {
			findings = append(findings, ResilienceFinding{
				Level:      "WARN",
				PolicyType: "MeshTimeout",
				Service:    service,
				Message:    "El servicio no está cubierto por ninguna política de timeouts.",
			})
		}
		if _, found := breakerCoveredServices[service]; !found {
			findings = append(findings, ResilienceFinding{
				Level:      "WARN",
				PolicyType: "MeshCircuitBreaker",
				Service:    service,
				Message:    "El servicio no está cubierto por ninguna política de circuit breaker.",
			})
		}
	}

	if len(findings) == 0 {
		findings = append(findings, ResilienceFinding{
			Level:   "INFO",
			Message: "Todos los servicios parecen tener políticas de resiliencia básicas aplicadas.",
		})
	}

	return &ValidationResult{
		Title:       "Análisis de Políticas de Resiliencia",
		GeneratedAt: time.Now(),
		Findings:    findings,
	}, nil
}

// getCoveredServices es una función helper para obtener los servicios cubiertos por un tipo de política.
func getCoveredServices(client dynamic.Interface, resourceName string, policyType string) (map[string]bool, error) {
	gvr := schema.GroupVersionResource{Group: "kuma.io", Version: "v1alpha1", Resource: resourceName}
	policies, err := client.Resource(gvr).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	coveredServices := make(map[string]bool)
	for _, policy := range policies.Items {
		to, found, _ := unstructured.NestedSlice(policy.Object, "spec", "to")
		if found {
			for _, toItem := range to {
				rule, _ := toItem.(map[string]interface{})
				targetRefName, _, _ := unstructured.NestedString(rule, "targetRef", "name")
				if targetRefName != "" {
					// Si la política apunta a '*', marcamos todos los servicios (aunque es una simplificación).
					// Por ahora, solo manejamos referencias directas a servicios.
					coveredServices[targetRefName] = true
				}
			}
		}
	}
	return coveredServices, nil
}

// getAllServices es una función helper para obtener un mapa de todos los servicios únicos del mesh.
func getAllServices(client dynamic.Interface) (map[string]bool, error) {
	dataplaneGVR := schema.GroupVersionResource{Group: "kuma.io", Version: "v1alpha1", Resource: "dataplanes"}
	dataplanes, err := client.Resource(dataplaneGVR).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error al listar Dataplanes para obtener servicios: %w", err)
	}

	allServices := make(map[string]bool)
	for _, dp := range dataplanes.Items {
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
	return allServices, nil
}
