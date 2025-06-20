// pkg/analysis/security.go
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

// AnalyzeMTLS revisa la configuración de mTLS en el Mesh y las políticas asociadas.
func AnalyzeMTLS(client dynamic.Interface) (*ValidationResult, error) {
	var findings []interface{}

	// Asumimos que el mesh a revisar se llama 'default'.
	// Una mejora futura sería permitir al usuario especificar el mesh.
	meshName := "default"
	meshGVR := schema.GroupVersionResource{Group: "kuma.io", Version: "v1alpha1", Resource: "meshes"}

	mesh, err := client.Resource(meshGVR).Get(context.TODO(), meshName, v1.GetOptions{})
	if err != nil {
		findings = append(findings, MTLSFinding{
			Level:    "ALERT",
			Message:  fmt.Sprintf("No se pudo obtener el Mesh '%s'. Error: %v", meshName, err),
			Resource: meshName,
		})
		return &ValidationResult{Title: "Análisis de Configuración mTLS", Findings: findings}, nil
	}

	// 1. Verificar si mTLS está habilitado en el Mesh
	enabledBackend, backendFound, _ := unstructured.NestedString(mesh.Object, "spec", "mtls", "enabledBackend")
	if !backendFound || enabledBackend == "" {
		findings = append(findings, MTLSFinding{
			Level:    "ALERT",
			Message:  "mTLS está DESACTIVADO para este mesh. El tráfico entre servicios no está cifrado.",
			Resource: meshName,
		})
	} else {
		findings = append(findings, MTLSFinding{
			Level:    "INFO",
			Message:  fmt.Sprintf("mTLS está ACTIVADO con el backend '%s'.", enabledBackend),
			Resource: meshName,
		})

		// 2. Verificar que el backend habilitado esté definido en la lista de backends
		backends, backendsFound, _ := unstructured.NestedSlice(mesh.Object, "spec", "mtls", "backends")
		if !backendsFound || len(backends) == 0 {
			findings = append(findings, MTLSFinding{
				Level:    "ALERT",
				Message:  fmt.Sprintf("El backend mTLS '%s' está habilitado, pero no se ha definido ninguna configuración de backends.", enabledBackend),
				Resource: meshName,
			})
		} else {
			isBackendDefined := false
			for _, backendItem := range backends {
				backendMap, _ := backendItem.(map[string]interface{})
				backendName, _, _ := unstructured.NestedString(backendMap, "name")
				if backendName == enabledBackend {
					isBackendDefined = true
					break
				}
			}
			if !isBackendDefined {
				findings = append(findings, MTLSFinding{
					Level:    "ALERT",
					Message:  fmt.Sprintf("El backend mTLS '%s' está habilitado, pero no se encuentra en la lista de backends definidos.", enabledBackend),
					Resource: meshName,
				})
			}
		}
	}

	// 3. Revisar MeshTrafficPermissions para ver si fuerzan mTLS
	mtpGVR := schema.GroupVersionResource{Group: "kuma.io", Version: "v1alpha1", Resource: "meshtrafficpermissions"}
	policies, err := client.Resource(mtpGVR).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error al listar MeshTrafficPermissions: %w", err)
	}

	for _, policy := range policies.Items {
		// --- INICIO DEL CÓDIGO CORREGIDO ---
		// La forma correcta de acceder a un campo dentro de un elemento de una lista (slice).

		// Primero, obtenemos la sección 'to' como una lista
		toSlice, found, _ := unstructured.NestedSlice(policy.Object, "spec", "to")
		if !found || len(toSlice) == 0 {
			continue // Si no hay sección 'to', no podemos analizarla.
		}

		// Accedemos al primer elemento de la lista
		firstToRule, ok := toSlice[0].(map[string]interface{})
		if !ok {
			continue // Si el formato no es el esperado
		}

		// Ahora, buscamos la acción dentro de ese primer elemento
		action, _, _ := unstructured.NestedString(firstToRule, "default", "action")
		// --- FIN DEL CÓDIGO CORREGIDO ---

		// Si mTLS está activo pero una política no lo fuerza, es una advertencia.
		if enabledBackend != "" && action != "" && action != "AllowWithMTLS" {
			findings = append(findings, MTLSFinding{
				Level:    "WARN",
				Message:  fmt.Sprintf("La política usa la acción '%s' en lugar de 'AllowWithMTLS', lo que podría permitir tráfico no cifrado.", action),
				Resource: policy.GetName(),
			})
		}
	}

	return &ValidationResult{
		Title:       "Análisis de Configuración mTLS",
		GeneratedAt: time.Now(),
		Findings:    findings,
	}, nil
}
