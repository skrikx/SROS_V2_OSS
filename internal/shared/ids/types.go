package ids

// Typed identifiers keep canonical contracts explicit and transport-neutral.
type (
	RunID            string
	TraceID          string
	SpanID           string
	EventID          string
	OperatorID       string
	TenantID         string
	WorkspaceID      string
	SessionID        string
	CheckpointID     string
	PolicyBundleID   string
	PolicyDecisionID string
	BranchID         string
	MemoryMutationID string
	ArtifactID       string
	EvidenceBundleID string
	ReceiptID        string
	ReleaseID        string
	RollbackID       string
)

const (
	DefaultTenantID    TenantID    = "local"
	DefaultWorkspaceID WorkspaceID = "default"
)

func (id RunID) String() string {
	return string(id)
}

func (id TraceID) String() string {
	return string(id)
}

func (id SpanID) String() string {
	return string(id)
}
