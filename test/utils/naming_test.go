package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestResourceName(t *testing.T) {
	tests := []struct {
		name         string
		suite        string
		testCase     string
		resourceType string
		expected     string
	}{
		{
			name:         "basic smoke test resource",
			suite:        "smoke",
			testCase:     "basic",
			resourceType: "model",
			expected:     "smoke-basic-model",
		},
		{
			name:         "saturation test resource",
			suite:        "saturation",
			testCase:     "scaleup",
			resourceType: "va",
			expected:     "saturation-scaleup-va",
		},
		{
			name:         "scale-to-zero test resource",
			suite:        "scale-to-zero",
			testCase:     "idle",
			resourceType: "deployment",
			expected:     "scale-to-zero-idle-deployment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TestResourceName(tt.suite, tt.testCase, tt.resourceType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTestResourceNameWithVariant(t *testing.T) {
	tests := []struct {
		name         string
		suite        string
		testCase     string
		resourceType string
		variant      string
		expected     string
	}{
		{
			name:         "variant A",
			suite:        "saturation",
			testCase:     "multi",
			resourceType: "model",
			variant:      "a",
			expected:     "saturation-multi-model-a",
		},
		{
			name:         "variant B",
			suite:        "saturation",
			testCase:     "multi",
			resourceType: "model",
			variant:      "b",
			expected:     "saturation-multi-model-b",
		},
		{
			name:         "cheap variant",
			suite:        "limiter",
			testCase:     "cost",
			resourceType: "va",
			variant:      "cheap",
			expected:     "limiter-cost-va-cheap",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TestResourceNameWithVariant(tt.suite, tt.testCase, tt.resourceType, tt.variant)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTestResourceNameWithSuffix(t *testing.T) {
	tests := []struct {
		name         string
		suite        string
		testCase     string
		resourceType string
		suffix       string
		expected     string
	}{
		{
			name:         "service suffix",
			suite:        "smoke",
			testCase:     "basic",
			resourceType: "model",
			suffix:       "service",
			expected:     "smoke-basic-model-service",
		},
		{
			name:         "decode suffix",
			suite:        "saturation",
			testCase:     "scaleup",
			resourceType: "model",
			suffix:       "decode",
			expected:     "saturation-scaleup-model-decode",
		},
		{
			name:         "monitor suffix",
			suite:        "pod-scraping",
			testCase:     "metrics",
			resourceType: "service",
			suffix:       "monitor",
			expected:     "pod-scraping-metrics-service-monitor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TestResourceNameWithSuffix(tt.suite, tt.testCase, tt.resourceType, tt.suffix)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsTestResource(t *testing.T) {
	tests := []struct {
		name     string
		resource string
		expected bool
	}{
		{
			name:     "smoke test resource",
			resource: "smoke-basic-model",
			expected: true,
		},
		{
			name:     "saturation test resource",
			resource: "saturation-scaleup-va",
			expected: true,
		},
		{
			name:     "scale-to-zero test resource",
			resource: "scale-to-zero-idle-deployment",
			expected: true,
		},
		{
			name:     "scale-from-zero test resource",
			resource: "scale-from-zero-startup-model",
			expected: true,
		},
		{
			name:     "parallel test resource",
			resource: "parallel-load-job-1",
			expected: true,
		},
		{
			name:     "limiter test resource",
			resource: "limiter-accelerator-va-a",
			expected: true,
		},
		{
			name:     "pod-scraping test resource",
			resource: "pod-scraping-metrics-service",
			expected: true,
		},
		{
			name:     "error-test resource",
			resource: "error-test-recovery-model",
			expected: true,
		},
		{
			name:     "target-condition resource",
			resource: "target-condition-ready-deployment",
			expected: true,
		},
		{
			name:     "generic test prefix",
			resource: "test-something",
			expected: true,
		},
		{
			name:     "production resource",
			resource: "production-model-service",
			expected: false,
		},
		{
			name:     "system resource",
			resource: "kube-system-pod",
			expected: false,
		},
		{
			name:     "user resource",
			resource: "my-application-deployment",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsTestResource(tt.resource)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewTestResourceNames(t *testing.T) {
	names := NewTestResourceNames("smoke", "basic")

	assert.Equal(t, "smoke-basic", names.Base)
	assert.Equal(t, "smoke-basic-decode", names.Deployment)
	assert.Equal(t, "smoke-basic-service", names.Service)
	assert.Equal(t, "smoke-basic-monitor", names.ServiceMonitor)
	assert.Equal(t, "smoke-basic-va", names.VA)
	assert.Equal(t, "smoke-basic-hpa", names.HPA)
	assert.Equal(t, "smoke-basic-so", names.ScaledObject)
	assert.Equal(t, "smoke-basic-pool", names.Pool)
	assert.Equal(t, "smoke-basic-load", names.LoadJob)
}

func TestNewTestResourceNamesWithVariant(t *testing.T) {
	namesA := NewTestResourceNamesWithVariant("saturation", "multi", "a")
	namesB := NewTestResourceNamesWithVariant("saturation", "multi", "b")

	// Verify variant A
	assert.Equal(t, "saturation-multi-a", namesA.Base)
	assert.Equal(t, "saturation-multi-a-decode", namesA.Deployment)
	assert.Equal(t, "saturation-multi-a-service", namesA.Service)
	assert.Equal(t, "saturation-multi-a-va", namesA.VA)

	// Verify variant B
	assert.Equal(t, "saturation-multi-b", namesB.Base)
	assert.Equal(t, "saturation-multi-b-decode", namesB.Deployment)
	assert.Equal(t, "saturation-multi-b-service", namesB.Service)
	assert.Equal(t, "saturation-multi-b-va", namesB.VA)

	// Verify variants are different
	assert.NotEqual(t, namesA.Base, namesB.Base)
	assert.NotEqual(t, namesA.Deployment, namesB.Deployment)
}

func TestResourceNameUniqueness(t *testing.T) {
	// Verify that different test cases produce unique names
	name1 := TestResourceName("smoke", "basic", "model")
	name2 := TestResourceName("smoke", "advanced", "model")
	name3 := TestResourceName("saturation", "basic", "model")

	assert.NotEqual(t, name1, name2, "Different test cases should produce different names")
	assert.NotEqual(t, name1, name3, "Different suites should produce different names")
	assert.NotEqual(t, name2, name3, "Different suites and test cases should produce different names")
}

func TestVariantNameUniqueness(t *testing.T) {
	// Verify that variants produce unique names
	nameA := TestResourceNameWithVariant("saturation", "multi", "model", "a")
	nameB := TestResourceNameWithVariant("saturation", "multi", "model", "b")
	nameC := TestResourceNameWithVariant("saturation", "multi", "model", "cheap")

	assert.NotEqual(t, nameA, nameB, "Different variants should produce different names")
	assert.NotEqual(t, nameA, nameC, "Different variants should produce different names")
	assert.NotEqual(t, nameB, nameC, "Different variants should produce different names")
}

// Made with Bob
