// pkg/analysis/types.go
package analysis

import "time"

// ValidationResult es una estructura genérica para contener los resultados de cualquier análisis.
type ValidationResult struct {
	Title       string        `json:"title"`
	GeneratedAt time.Time     `json:"generatedAt"`
	Findings    []interface{} `json:"findings"` // Usamos interface{} para poder guardar cualquier tipo de hallazgo.
}

// DataplaneStatus contiene el estado de salud de un único Dataplane.
type DataplaneStatus struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Status    string `json:"status"`
	Details   string `json:"details"`
}

// SummaryStatus contiene los datos para el resumen general del mesh.
type SummaryStatus struct {
	TotalMeshes        int `json:"totalMeshes"`
	TotalDataplanes    int `json:"totalDataplanes"`
	OnlineDataplanes   int `json:"onlineDataplanes"`
	OfflineDataplanes  int `json:"offlineDataplanes"`
	DegradedDataplanes int `json:"degradedDataplanes"`
	InfoDataplanes     int `json:"infoDataplanes"`
	TotalPolicies      int `json:"totalPolicies"`
}

// PolicyFinding representa un hallazgo (problema o información) sobre una política.
type PolicyFinding struct {
	Level    string `json:"level"` // "ALERT", "WARN", "INFO"
	Message  string `json:"message"`
	Resource string `json:"resource"` // El nombre del recurso asociado (servicio o política)
}

// MTLSFinding representa un hallazgo sobre la configuración de mTLS.
type MTLSFinding struct {
	Level    string `json:"level"` // "INFO", "WARN", "ALERT"
	Message  string `json:"message"`
	Resource string `json:"resource"` // El Mesh o la política específica
}

// ResilienceFinding representa un hallazgo sobre las políticas de resiliencia.
type ResilienceFinding struct {
	Level      string `json:"level"`      // "WARN", "INFO"
	PolicyType string `json:"policyType"` // "MeshRetry", "MeshTimeout", etc.
	Service    string `json:"service"`
	Message    string `json:"message"`
}

// ObservabilityFinding representa un hallazgo sobre las políticas de observabilidad.
type ObservabilityFinding struct {
	Level      string `json:"level"`      // "WARN", "INFO"
	PolicyType string `json:"policyType"` // "MeshLog", "MeshMetric", "MeshTrace"
	Resource   string `json:"resource"`   // El nombre de la política o "Global"
	Message    string `json:"message"`
}
