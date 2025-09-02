package testdata

import (
	"time"
)

// +flags-gen
// OperatorConfig defines configuration options for the operator
type OperatorConfig struct {

	// Controllers is a list of controllers to enable. '*' enables all on-by-default controllers,
	// 'foo' enables the controller named 'foo', '-foo' disables the controller named 'foo'.
	// +optional
	// +listType=set
	Controllers []string `json:"controllers,omitempty" yaml:"controllers,omitempty" default:"*"`

	// ProbeAddr is the address the probe endpoint binds to.
	// +optional
	ProbeAddr string `json:"probeAddr,omitempty" yaml:"probeAddr,omitempty" default:":8080"`

	// MetricsAddr is the address the metrics endpoint binds to.
	// +optional
	MetricsAddr string `json:"metricsAddr,omitempty" yaml:"metricsAddr,omitempty" default:":8443"`

	// ProbeHealthEndpoint is the endpoint for the health probe.
	// +optional
	ProbeHealthEndpoint string `json:"probeHealthEndpoint,omitempty" yaml:"probeHealthEndpoint,omitempty" default:"healthz"`

	// ProbeReadyEndpoint is the endpoint for the ready probe.
	// +optional
	ProbeReadyEndpoint string `json:"probeReadyEndpoint,omitempty" yaml:"probeReadyEndpoint,omitempty" default:"readyz"`

	// EnableLeaderElection enables leader election for controller manager.
	// Enabling this will ensure there is only one active controller manager.
	// +optional
	EnableLeaderElection bool `json:"enableLeaderElection,omitempty" yaml:"enableLeaderElection,omitempty"`

	// ZapDevMode enables development mode for zap logger.
	// Enabling this will use human-readable output instead of structured JSON.
	// +optional
	ZapDevMode bool `json:"zapDevMode,omitempty" yaml:"zapDevMode,omitempty"`

	// V is the log level for V logs.
	// +optional
	V int `json:"v,omitempty" yaml:"v,omitempty"`

	// RequiredCRDs is a list of CRDs that must be present before starting the controller manager.
	// Format: group/version/kind. Example: eventing.knative.dev/v1/Broker
	// +optional
	// +listType=set
	RequiredCRDs []string `json:"requiredCRDs,omitempty" yaml:"requiredCRDs,omitempty" default:"eventing.knative.dev/v1/Broker,eventing.knative.dev/v1/Trigger,serving.knative.dev/v1/Service,sources.knative.dev/v1/SinkBinding"`

	// RequiredCRDsGracePeriod is the grace period for the required CRDs to be present
	// before starting the controller manager.
	// +optional
	RequiredCRDsGracePeriod time.Duration `json:"requiredCRDsGracePeriod,omitempty" yaml:"requiredCRDsGracePeriod,omitempty" default:"30s"`

	// RuntimeConfigMapName is the name of the runtime config map.
	// +optional
	RuntimeConfigMapName string `json:"runtimeConfigMapName,omitempty" yaml:"runtimeConfigMapName,omitempty" default:"runtime-configmap"`

	// RuntimeConfigMapNamespace is the namespace of the runtime config map.
	// +optional
	RuntimeConfigMapNamespace string `json:"runtimeConfigMapNamespace,omitempty" yaml:"runtimeConfigMapNamespace,omitempty"`

	// RuntimeConfigKey is the key of the runtime config in the configmap.
	// +optional
	RuntimeConfigKey string `json:"runtimeConfigKey,omitempty" yaml:"runtimeConfigKey,omitempty" default:"runtime-config.yaml"`
}
