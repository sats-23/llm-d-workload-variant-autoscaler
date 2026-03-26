# Justified Ensure* Usage in E2E Tests

This document catalogs the remaining `Ensure*` function usage in the e2e test suite and provides justification for each case.

## Overview

As part of the test isolation refactoring effort, we've established that `Create*` functions with unique resource names should be preferred over `Ensure*` functions. However, certain scenarios justify the use of `Ensure*` semantics.

## Justification Criteria

`Ensure*` usage is justified when:

1. **Testing Recovery Behavior**: The test explicitly validates controller recovery when resources are deleted/recreated
2. **Shared Test Infrastructure**: Resources are truly shared across all tests (e.g., namespace-level configuration)
3. **Development/Debugging**: Temporary usage during development with clear TODO comments
4. **Idempotent Retries**: Test framework retries require idempotent setup

## Current Justified Usage

### 1. Error Recovery Tests (smoke_test.go)

**Location**: `test/e2e/smoke_test.go:1068-1070`

```go
// JUSTIFIED: Testing controller recovery behavior when deployment is recreated
By("Recreating the deployment")
err = fixtures.EnsureModelService(ctx, k8sClient, cfg.LLMDNamespace, 
    errorTestModelServiceName, errorTestPoolName, cfg.ModelID, 
    cfg.UseSimulator, cfg.MaxNumSeqs)
```

**Justification**: This test specifically validates that the WVA controller can recover when a deployment is deleted and recreated. The `Ensure*` pattern is intentional here to simulate real-world scenarios where deployments might be recreated by operators or other controllers.

**Alternative Considered**: Using `Delete` + `Create` explicitly, but `Ensure*` better represents the real-world scenario being tested.

**Status**: ✅ Justified - Keep with documentation

---

### 2. Remaining Unjustified Usage

The following `Ensure*` calls should be migrated to `Create*` with unique names:

#### smoke_test.go

| Line | Function | Context | Migration Priority |
|------|----------|---------|-------------------|
| 135 | `EnsureModelService` | Basic VA lifecycle setup | High |
| 151 | `EnsureService` | Basic VA lifecycle setup | High |
| 168 | `EnsureServiceMonitor` | Basic VA lifecycle setup | High |
| 197 | `EnsureVariantAutoscalingWithDefaults` | Basic VA lifecycle setup | High |
| 210 | `EnsureScaledObject` | Basic VA lifecycle setup | High |
| 213 | `EnsureHPA` | Basic VA lifecycle setup | High |
| 590 | `EnsureBurstLoadJob` | Scale-up test | High |
| 974 | `EnsureModelService` | Error test setup | Medium |
| 997 | `EnsureVariantAutoscalingWithDefaults` | Error test setup | Medium |

**Migration Plan**: 
- Use `utils.NewTestResourceNames("smoke", "basic")` for basic lifecycle tests
- Use `utils.NewTestResourceNames("smoke", "scaleup")` for scale-up tests
- Use `utils.NewTestResourceNames("smoke", "error")` for error handling tests

#### saturation_test.go

| Line | Function | Context | Migration Priority |
|------|----------|---------|-------------------|
| 42 | `EnsureModelService` | Single variant test | High |
| 58 | `EnsureService` | Single variant test | High |
| 74 | `EnsureServiceMonitor` | Single variant test | High |
| 129 | `EnsureVariantAutoscaling` | Single variant test | High |
| 145 | `EnsureScaledObject` | Single variant test | High |
| 148 | `EnsureHPA` | Single variant test | High |
| 341-349 | Multiple `Ensure*` | Multi-variant test A | High |
| 352-360 | Multiple `Ensure*` | Multi-variant test B | High |
| 375-380 | `EnsureVariantAutoscaling` | Multi-variant VAs | High |
| 387-395 | `EnsureScaledObject`/`EnsureHPA` | Multi-variant scalers | High |
| 451-456 | `EnsureBurstLoadJob` | Load generation | Medium |

**Migration Plan**:
- Use `utils.NewTestResourceNames("saturation", "single")` for single variant
- Use `utils.NewTestResourceNamesWithVariant("saturation", "multi", "a")` and `"b"` for multi-variant

#### scale_to_zero_test.go

| Line | Function | Context | Migration Priority |
|------|----------|---------|-------------------|
| 47 | `EnsureModelService` | Scale-to-zero test | High |
| 63 | `EnsureService` | Scale-to-zero test | High |
| 79 | `EnsureServiceMonitor` | Scale-to-zero test | High |
| 108 | `EnsureVariantAutoscaling` | Scale-to-zero test | High |
| 119 | `EnsureScaledObject` | Scale-to-zero test | High |
| 122 | `EnsureHPA` | Scale-to-zero test | High |
| 352 | `EnsureModelService` | Disabled test | Medium |
| 357 | `EnsureService` | Disabled test | Medium |
| 362 | `EnsureServiceMonitor` | Disabled test | Medium |
| 375 | `EnsureVariantAutoscaling` | Disabled test | Medium |
| 385 | `EnsureScaledObject` | Disabled test | Medium |
| 388 | `EnsureHPA` | Disabled test | Medium |

**Migration Plan**:
- Use `utils.NewTestResourceNames("scale-to-zero", "idle")` for main test
- Use `utils.NewTestResourceNames("scale-to-zero", "disabled")` for disabled test

#### scale_from_zero_test.go

| Line | Function | Context | Migration Priority |
|------|----------|---------|-------------------|
| 94 | `EnsureModelService` | Scale-from-zero test | High |
| 108 | `EnsureService` | Scale-from-zero test | High |
| 112 | `EnsureServiceMonitor` | Scale-from-zero test | High |
| 141 | `EnsureVariantAutoscaling` | Scale-from-zero test | High |
| 152 | `EnsureScaledObject` | Scale-from-zero test | High |
| 155 | `EnsureHPA` | Scale-from-zero test | High |

**Migration Plan**:
- Use `utils.NewTestResourceNames("scale-from-zero", "startup")` 

#### parallel_load_scaleup_test.go

| Line | Function | Context | Migration Priority |
|------|----------|---------|-------------------|
| 92 | `EnsureModelService` | Parallel load test | High |
| 108 | `EnsureService` | Parallel load test | High |
| 124 | `EnsureServiceMonitor` | Parallel load test | High |
| 153 | `EnsureVariantAutoscalingWithDefaults` | Parallel load test | High |
| 186 | `EnsureScaledObject` | Parallel load test | High |
| 192 | `EnsureHPA` | Parallel load test | High |
| 321 | `EnsureParallelLoadJobs` | Load generation | Medium |

**Migration Plan**:
- Use `utils.NewTestResourceNames("parallel", "scaleup")`

#### limiter_test.go

| Line | Function | Context | Migration Priority |
|------|----------|---------|-------------------|
| 43-51 | Multiple `Ensure*` | Limiter test variant A | High |
| 72-80 | Multiple `Ensure*` | Limiter test variant B | High |
| 114-123 | `EnsureVariantAutoscaling` | Limiter VAs | High |
| 133-141 | `EnsureScaledObject`/`EnsureHPA` | Limiter scalers | High |

**Migration Plan**:
- Use `utils.NewTestResourceNamesWithVariant("limiter", "accelerator", "nvidia")`
- Use `utils.NewTestResourceNamesWithVariant("limiter", "accelerator", "amd")`

#### pod_scraping_test.go

| Line | Function | Context | Migration Priority |
|------|----------|---------|-------------------|
| 34 | `EnsureModelService` | Pod scraping test | Low |
| 39 | `EnsureService` | Pod scraping test | Low |

**Migration Plan**:
- Use `utils.NewTestResourceNames("pod-scraping", "metrics")`

## Migration Priority

### High Priority (Core Test Functionality)
- smoke_test.go: Basic VA lifecycle (lines 135-213)
- saturation_test.go: All tests
- scale_to_zero_test.go: Main test
- scale_from_zero_test.go: All tests
- parallel_load_scaleup_test.go: All tests
- limiter_test.go: All tests

### Medium Priority (Secondary Tests)
- smoke_test.go: Error handling tests
- scale_to_zero_test.go: Disabled test
- Load generation jobs

### Low Priority (Infrastructure Tests)
- pod_scraping_test.go

## Migration Tracking

### Completed
- ✅ Infrastructure and utilities created
- ✅ Documentation written
- ✅ Fixture functions documented

### In Progress
- 🔄 Migration guide created
- 🔄 Justified usage documented

### Pending
- ⏳ Actual test file migrations (54 Ensure* calls across 7 files)

## Adding New Justified Usage

If you need to add new `Ensure*` usage, follow this process:

1. **Document in Code**: Add a comment explaining why `Ensure*` is needed
   ```go
   // JUSTIFIED ENSURE* USAGE: [Reason]
   // Alternative considered: [What was considered and why it doesn't work]
   err := fixtures.EnsureModelService(...)
   ```

2. **Update This Document**: Add an entry to the "Current Justified Usage" section

3. **Code Review**: Ensure reviewers understand and approve the justification

4. **Periodic Review**: Revisit justified usage quarterly to see if it can be eliminated

## See Also

- [Migration Guide](./e2e-test-isolation-migration-guide.md)
- [Refactoring Plan](./e2e-test-isolation-refactoring.md)
- [Implementation Summary](./e2e-test-isolation-implementation-summary.md)