package memory

type Scope string

const (
	ScopeSession   Scope = "session"
	ScopeWorkspace Scope = "workspace"
	ScopeRun       Scope = "run"
	ScopeGlobal    Scope = "global"
)
