package policy

type TrustBoundary string

const (
	TrustBoundaryLocalProcess TrustBoundary = "local_process"
	TrustBoundaryLocalFS      TrustBoundary = "local_fs"
	TrustBoundaryExternalNet  TrustBoundary = "external_net"
)
