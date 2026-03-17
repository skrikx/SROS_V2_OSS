package sr8

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"srosv2/contracts/evidence"
	"srosv2/contracts/runcontract"
)

type Options struct {
	Now func() time.Time
}

type Compiler struct {
	now func() time.Time
}

func NewCompiler(opts Options) *Compiler {
	now := opts.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return &Compiler{now: now}
}

func (c *Compiler) Compile(_ context.Context, req CompileRequest) (CompileResult, error) {
	now := c.now().UTC()

	parsed, err := ParseRequest(req)
	if err != nil {
		return CompileResult{}, err
	}
	normalized := Normalize(parsed, now)
	classification := Classify(normalized)
	topology := SelectTopology(normalized, classification)

	if err := ValidateCompileInput(normalized, classification, topology); err != nil {
		return CompileResult{}, err
	}

	contract := AssembleRunContract(normalized, classification, topology, now)
	outputDir := filepath.Join(defaultArtifactRoot(req.ArtifactRoot), "compile", string(contract.RunID))
	runContractPath := filepath.Join(outputDir, "run_contract.json")
	receiptPath := filepath.Join(outputDir, "compile_receipt.json")
	srxmlPath := filepath.Join(outputDir, "run_contract.srxml")

	contract.ArtifactRefs = plannedArtifactRefs(runContractPath, receiptPath, srxmlPath, req.EmitSRXML)
	if err := ValidateRunContract(contract); err != nil {
		return CompileResult{}, err
	}

	srxmlText, err := RenderSRXML(contract)
	if err != nil {
		return CompileResult{}, err
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return CompileResult{}, fmt.Errorf("create compile output directory: %w", err)
	}
	if err := writeJSONFile(runContractPath, contract); err != nil {
		return CompileResult{}, err
	}
	if req.EmitSRXML {
		if err := os.WriteFile(srxmlPath, []byte(srxmlText), 0o644); err != nil {
			return CompileResult{}, fmt.Errorf("write srxml output: %w", err)
		}
	}

	bundleArtifacts, err := compileBundleArtifacts(runContractPath, srxmlPath, req.EmitSRXML)
	if err != nil {
		return CompileResult{}, err
	}
	receipt, err := BuildCompileReceipt(normalized.CompileRequestID, classification, topology, contract, bundleArtifacts, now)
	if err != nil {
		return CompileResult{}, err
	}
	if err := writeJSONFile(receiptPath, receipt); err != nil {
		return CompileResult{}, err
	}

	summary := fmt.Sprintf("compiled %s as %s (%s)", contract.RunID, classification.Domain, classification.Risk)
	return CompileResult{
		Accepted:       true,
		Summary:        summary,
		Normalized:     normalized,
		Classification: classification,
		Topology:       topology,
		RunContract:    contract,
		SRXML:          srxmlText,
		Receipt:        receipt,
		Artifacts:      contract.ArtifactRefs,
		OutputDir:      outputDir,
	}, nil
}

func defaultArtifactRoot(v string) string {
	if v == "" {
		return filepath.Clean(filepath.Join(".", "artifacts"))
	}
	return filepath.Clean(v)
}

func plannedArtifactRefs(runContractPath, receiptPath, srxmlPath string, emitSRXML bool) []runcontract.ArtifactReference {
	refs := []runcontract.ArtifactReference{
		{Kind: "run_contract", ArtifactID: "artifact_run_contract", URI: filepath.Clean(runContractPath)},
		{Kind: "compile_receipt", ArtifactID: "artifact_compile_receipt", URI: filepath.Clean(receiptPath)},
	}
	if emitSRXML {
		refs = append(refs, runcontract.ArtifactReference{Kind: "srxml", ArtifactID: "artifact_srxml", URI: filepath.Clean(srxmlPath)})
	}
	return refs
}

func compileBundleArtifacts(runContractPath, srxmlPath string, emitSRXML bool) ([]evidence.ArtifactRef, error) {
	refs := make([]evidence.ArtifactRef, 0, 2)
	runRef, err := ArtifactRefFromFile(runContractPath)
	if err != nil {
		return nil, fmt.Errorf("build run contract artifact ref: %w", err)
	}
	refs = append(refs, runRef)

	if emitSRXML {
		xref, err := ArtifactRefFromFile(srxmlPath)
		if err != nil {
			return nil, fmt.Errorf("build srxml artifact ref: %w", err)
		}
		refs = append(refs, xref)
	}

	return refs, nil
}

func writeJSONFile(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json %s: %w", path, err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("write file %s: %w", path, err)
	}
	return nil
}
