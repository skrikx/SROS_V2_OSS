package mcpclient

import (
	ctools "srosv2/contracts/tools"
)

type RemoteCapability struct {
	Name         string            `json:"name"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	Domain       string            `json:"domain"`
	Endpoint     string            `json:"endpoint"`
	AuthType     string            `json:"auth_type"`
	Capabilities []string          `json:"capabilities"`
	Metadata     map[string]string `json:"metadata"`
}

func Normalize(remote RemoteCapability) ctools.Manifest {
	return ctools.Manifest{
		ManifestVersion:      "v2.0",
		Name:                 remote.Name,
		Title:                remote.Title,
		Description:          remote.Description,
		Version:              "ingested",
		Class:                "mcp.ingress",
		Domain:               remote.Domain,
		PolicyClass:          "connector.governed",
		Status:               ctools.StateValidated,
		TrustBoundary:        "external_net",
		SandboxProfiles:      []string{"network_restricted"},
		AuthType:             remote.AuthType,
		AuthEnvelopeRequired: remote.AuthType != "",
		MCPIngressCapable:    true,
		RemoteCapable:        true,
		Capabilities:         remote.Capabilities,
		Metadata:             remote.Metadata,
		ConnectorRef:         remote.Endpoint,
	}
}
