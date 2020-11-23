package hit

// This file exposes functionality that should only be accessible during tests

// cast to uint8 to avoid linter problems in step.go
const CleanStep = uint8(cleanStep)

type CallPath = callPath

var NewCallPath = newCallPath
