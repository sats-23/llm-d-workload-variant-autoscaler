package utils

import (
	"fmt"
	"strings"
)

// TestResourceName generates a unique resource name for a test.
// This ensures test isolation by preventing resource name collisions.
//
// Parameters:
//   - suite: The test suite name (e.g., "smoke", "saturation", "scale-to-zero")
//   - testCase: The specific test case name (e.g., "basic", "multi-variant", "idle")
//   - resourceType: The type of resource (e.g., "model", "service", "va", "hpa")
//
// Example:
//
//	name := TestResourceName("smoke", "basic", "model")
//	// Returns: "smoke-basic-model"
func TestResourceName(suite, testCase, resourceType string) string {
	return fmt.Sprintf("%s-%s-%s", suite, testCase, resourceType)
}

// TestResourceNameWithVariant generates a unique resource name with a variant suffix.
// Use this when a test creates multiple instances of the same resource type.
//
// Parameters:
//   - suite: The test suite name
//   - testCase: The specific test case name
//   - resourceType: The type of resource
//   - variant: The variant identifier (e.g., "a", "b", "cheap", "expensive")
//
// Example:
//
//	nameA := TestResourceNameWithVariant("saturation", "multi", "model", "a")
//	// Returns: "saturation-multi-model-a"
//	nameB := TestResourceNameWithVariant("saturation", "multi", "model", "b")
//	// Returns: "saturation-multi-model-b"
func TestResourceNameWithVariant(suite, testCase, resourceType, variant string) string {
	return fmt.Sprintf("%s-%s-%s-%s", suite, testCase, resourceType, variant)
}

// TestResourceNameWithSuffix generates a unique resource name with a custom suffix.
// This is useful for derived resource names (e.g., "model-service", "model-decode").
//
// Parameters:
//   - suite: The test suite name
//   - testCase: The specific test case name
//   - resourceType: The type of resource
//   - suffix: The suffix to append (e.g., "service", "decode", "monitor")
//
// Example:
//
//	baseName := TestResourceName("smoke", "basic", "model")
//	serviceName := TestResourceNameWithSuffix("smoke", "basic", "model", "service")
//	// Returns: "smoke-basic-model-service"
func TestResourceNameWithSuffix(suite, testCase, resourceType, suffix string) string {
	return fmt.Sprintf("%s-%s-%s-%s", suite, testCase, resourceType, suffix)
}

// IsTestResource checks if a resource name matches test resource naming patterns.
// This is used by cleanup functions to identify test resources.
//
// A resource is considered a test resource if its name starts with any of:
//   - "test-"
//   - "smoke-"
//   - "saturation-"
//   - "scale-to-zero-"
//   - "scale-from-zero-"
//   - "parallel-"
//   - "limiter-"
//   - "pod-scraping-"
//   - "error-test-"
//   - "target-condition-"
func IsTestResource(name string) bool {
	testPrefixes := []string{
		"test-",
		"smoke-",
		"saturation-",
		"scale-to-zero-",
		"scale-from-zero-",
		"parallel-",
		"limiter-",
		"pod-scraping-",
		"error-test-",
		"target-condition-",
	}

	for _, prefix := range testPrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

// TestResourceNames is a helper struct that holds all related resource names
// for a test case. This makes it easier to manage multiple related resources.
type TestResourceNames struct {
	// Base name for the test resources
	Base string

	// Model service deployment name (typically base + "-decode")
	Deployment string

	// Service name (typically base + "-service")
	Service string

	// ServiceMonitor name (typically base + "-monitor")
	ServiceMonitor string

	// VariantAutoscaling name (typically base + "-va")
	VA string

	// HPA name (typically base + "-hpa")
	HPA string

	// ScaledObject name (typically base + "-so")
	ScaledObject string

	// InferencePool name (typically base + "-pool")
	Pool string

	// Load job name (typically base + "-load")
	LoadJob string
}

// NewTestResourceNames creates a TestResourceNames struct with standard naming conventions.
//
// Parameters:
//   - suite: The test suite name
//   - testCase: The specific test case name
//
// Example:
//
//	names := NewTestResourceNames("smoke", "basic")
//	// names.Base = "smoke-basic"
//	// names.Deployment = "smoke-basic-decode"
//	// names.Service = "smoke-basic-service"
//	// names.VA = "smoke-basic-va"
//	// etc.
func NewTestResourceNames(suite, testCase string) TestResourceNames {
	base := fmt.Sprintf("%s-%s", suite, testCase)
	return TestResourceNames{
		Base:           base,
		Deployment:     base + "-decode",
		Service:        base + "-service",
		ServiceMonitor: base + "-monitor",
		VA:             base + "-va",
		HPA:            base + "-hpa",
		ScaledObject:   base + "-so",
		Pool:           base + "-pool",
		LoadJob:        base + "-load",
	}
}

// NewTestResourceNamesWithVariant creates a TestResourceNames struct for a specific variant.
// Use this when a test creates multiple sets of resources (e.g., variant A and variant B).
//
// Parameters:
//   - suite: The test suite name
//   - testCase: The specific test case name
//   - variant: The variant identifier
//
// Example:
//
//	namesA := NewTestResourceNamesWithVariant("saturation", "multi", "a")
//	// namesA.Base = "saturation-multi-a"
//	// namesA.Deployment = "saturation-multi-a-decode"
//	// etc.
//
//	namesB := NewTestResourceNamesWithVariant("saturation", "multi", "b")
//	// namesB.Base = "saturation-multi-b"
//	// namesB.Deployment = "saturation-multi-b-decode"
//	// etc.
func NewTestResourceNamesWithVariant(suite, testCase, variant string) TestResourceNames {
	base := fmt.Sprintf("%s-%s-%s", suite, testCase, variant)
	return TestResourceNames{
		Base:           base,
		Deployment:     base + "-decode",
		Service:        base + "-service",
		ServiceMonitor: base + "-monitor",
		VA:             base + "-va",
		HPA:            base + "-hpa",
		ScaledObject:   base + "-so",
		Pool:           base + "-pool",
		LoadJob:        base + "-load",
	}
}
