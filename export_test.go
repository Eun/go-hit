package hit

// This file exposes functionality that should only be accessible during tests

// CleanStep is a step that runs during the clean step phase, cast to uint8 to avoid linter problems in step.go.
const CleanStep = uint8(cleanStep)

// CallPath represents the internal callPath.
type CallPath = callPath

// NewCallPath creates a new callPath.
var NewCallPath = newCallPath
