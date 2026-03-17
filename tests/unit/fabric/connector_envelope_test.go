package fabric_test

import (
	"testing"

	"srosv2/internal/fabric/connectors"
)

func TestConnectorEnvelopeRedactsSecretMaterial(t *testing.T) {
	env, err := connectors.LoadEnvelope("../../../examples/connectors/local_secret_envelope.json")
	if err != nil {
		t.Fatalf("load envelope: %v", err)
	}
	redacted := connectors.RedactEnvelope(env)
	if redacted["secret_material"] != "[REDACTED]" {
		t.Fatalf("expected redacted secret material")
	}
}
