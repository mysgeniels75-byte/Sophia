// Package validation provides constraint validation for OSE Advisory Service.
//
// This implements the "validation at the boundary" pattern where all defensive
// programming is concentrated at the API Gateway, creating a trust boundary
// between untrusted external inputs and trusted internal systems.
package validation

import (
	"fmt"
	"regexp"
	"strings"

	pb "github.com/your-org/ose-api/api/proto/advisory/v1"
)

// Service name must be DNS-compatible: lowercase, hyphens, 3-63 characters
var serviceNameRegex = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// ValidateServiceConstraints validates all aspects of service constraints.
//
// This function implements the mathematical invariants specified in the
// architectural documentation:
//
//   1. Structural Validity (service name format)
//   2. Numerical Constraints (positive values, reasonable ranges)
//   3. Enumeration Constraints (valid enum values)
//   4. Cross-Field Constraints (interdependent validation rules)
//   5. Integration Constraints (valid integration combinations)
//
// Returns:
//   nil if all constraints are valid
//   ValidationError describing the first violation encountered
func ValidateServiceConstraints(c *pb.ServiceConstraints) error {
	if c == nil {
		return &ValidationError{
			Field:   "constraints",
			Message: "constraints cannot be nil",
		}
	}

	// ═══════════════════════════════════════════════════════════════════════
	// INVARIANT 1: STRUCTURAL VALIDITY
	// ═══════════════════════════════════════════════════════════════════════

	// Service name validation
	if c.ServiceName == "" {
		return &ValidationError{
			Field:   "service_name",
			Message: "service name is required",
		}
	}

	if len(c.ServiceName) < 3 || len(c.ServiceName) > 63 {
		return &ValidationError{
			Field:   "service_name",
			Message: fmt.Sprintf("service name must be 3-63 characters, got %d", len(c.ServiceName)),
			Suggestion: "Use a concise, descriptive name like 'inventory-manager' or 'payment-processor'",
		}
	}

	if !serviceNameRegex.MatchString(c.ServiceName) {
		return &ValidationError{
			Field:   "service_name",
			Message: "service name must be lowercase letters, numbers, and hyphens only (DNS-compatible)",
			Suggestion: fmt.Sprintf("Try: %s", strings.ToLower(strings.ReplaceAll(c.ServiceName, "_", "-"))),
		}
	}

	// ═══════════════════════════════════════════════════════════════════════
	// INVARIANT 2: NUMERICAL CONSTRAINTS
	// ═══════════════════════════════════════════════════════════════════════

	// Throughput validation
	if c.ThroughputTps <= 0 {
		return &ValidationError{
			Field:   "throughput_tps",
			Message: "throughput must be positive",
			Suggestion: "Specify expected peak requests/messages per second (e.g., 1000 for typical API)",
		}
	}

	if c.ThroughputTps > 1_000_000 {
		return &ValidationError{
			Field:   "throughput_tps",
			Message: fmt.Sprintf("throughput %d TPS exceeds reasonable maximum (1M TPS)", c.ThroughputTps),
			Suggestion: "If you truly need >1M TPS, contact the OSE team for specialized architecture guidance",
		}
	}

	// Latency validation
	if c.LatencyP99Ms <= 0 {
		return &ValidationError{
			Field:   "latency_p99_ms",
			Message: "latency target must be positive",
			Suggestion: "Specify p99 latency target in milliseconds (e.g., 100 for typical API)",
		}
	}

	if c.LatencyP99Ms > 60_000 {
		return &ValidationError{
			Field:   "latency_p99_ms",
			Message: fmt.Sprintf("latency target %dms exceeds 60 seconds", c.LatencyP99Ms),
			Suggestion: "For batch jobs with >60s latency, consider ServiceType = BACKGROUND_WORKER",
		}
	}

	// ═══════════════════════════════════════════════════════════════════════
	// INVARIANT 3: ENUMERATION CONSTRAINTS
	// ═══════════════════════════════════════════════════════════════════════

	// Service type validation
	if c.ServiceType == pb.ServiceType_SERVICE_TYPE_UNSPECIFIED {
		return &ValidationError{
			Field:   "service_type",
			Message: "service type must be specified",
			Suggestion: "Choose: API, EVENT_PROCESSOR, BACKGROUND_WORKER, or STREAM_PROCESSOR",
		}
	}

	// Consistency model validation
	if c.ConsistencyModel == pb.ConsistencyModel_CONSISTENCY_MODEL_UNSPECIFIED {
		return &ValidationError{
			Field:   "consistency_model",
			Message: "consistency model must be specified",
			Suggestion: "Choose STRONG (ACID transactions) or EVENTUAL (BASE properties)",
		}
	}

	// Deployment target validation
	if c.DeploymentTarget == pb.DeploymentTarget_DEPLOYMENT_TARGET_UNSPECIFIED {
		return &ValidationError{
			Field:   "deployment_target",
			Message: "deployment target must be specified",
			Suggestion: "Choose: KUBERNETES, ECS, or LAMBDA",
		}
	}

	// ═══════════════════════════════════════════════════════════════════════
	// INVARIANT 4: CROSS-FIELD CONSTRAINTS
	// ═══════════════════════════════════════════════════════════════════════

	// Strong consistency implies minimum latency overhead
	if c.ConsistencyModel == pb.ConsistencyModel_CONSISTENCY_MODEL_STRONG && c.LatencyP99Ms < 50 {
		return &ValidationError{
			Field:   "latency_p99_ms",
			Message: "strong consistency requires minimum 50ms latency (ACID coordination overhead)",
			Suggestion: "Either increase latency target to ≥50ms or use eventual consistency",
		}
	}

	// High throughput requires appropriate deployment target
	if c.ThroughputTps > 10_000 && c.DeploymentTarget == pb.DeploymentTarget_DEPLOYMENT_TARGET_LAMBDA {
		return &ValidationError{
			Field:   "deployment_target",
			Message: "Lambda is not suitable for >10K TPS sustained throughput",
			Suggestion: "Use KUBERNETES for high-throughput services",
		}
	}

	// ═══════════════════════════════════════════════════════════════════════
	// INVARIANT 5: INTEGRATION CONSTRAINTS
	// ═══════════════════════════════════════════════════════════════════════

	if len(c.Integrations) > 10 {
		return &ValidationError{
			Field:   "integrations",
			Message: fmt.Sprintf("service declares %d integrations (max 10)", len(c.Integrations)),
			Suggestion: "Services with >10 dependencies are likely violating single responsibility principle",
		}
	}

	// Validate each integration type is valid
	for _, integration := range c.Integrations {
		if integration == pb.IntegrationType_INTEGRATION_TYPE_UNSPECIFIED {
			return &ValidationError{
				Field:   "integrations",
				Message: "integration type UNSPECIFIED is not allowed",
				Suggestion: "Specify concrete integration types (KAFKA, POSTGRESQL, etc.)",
			}
		}
	}

	// All validations passed
	return nil
}

// ValidationError represents a constraint validation failure with context.
type ValidationError struct {
	Field      string // Which field failed validation
	Message    string // What went wrong
	Suggestion string // How to fix it (constructive validation)
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	if e.Suggestion != "" {
		return fmt.Sprintf("%s: %s. Suggestion: %s", e.Field, e.Message, e.Suggestion)
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}
