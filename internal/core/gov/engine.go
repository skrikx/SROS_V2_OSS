package gov

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"srosv2/contracts/policy"
)

type Options struct {
	BundlePath   string
	ArtifactRoot string
	Now          func() time.Time
	Bundle       *policy.Bundle
}

type Engine struct {
	bundle    policy.Bundle
	auditRoot string
	now       func() time.Time
}

func NewEngine(opts Options) (*Engine, error) {
	bundle := policy.Bundle{}
	switch {
	case opts.Bundle != nil:
		bundle = *opts.Bundle
	case strings.TrimSpace(opts.BundlePath) != "":
		loaded, err := LoadBundle(opts.BundlePath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				bundle = DefaultBundle()
				break
			}
			return nil, err
		}
		bundle = loaded
	default:
		bundle = DefaultBundle()
	}
	now := opts.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return &Engine{
		bundle:    bundle,
		auditRoot: filepath.Join(opts.ArtifactRoot, "gov"),
		now:       now,
	}, nil
}

func DefaultBundle() policy.Bundle {
	return policy.Bundle{
		BundleID:              "pb_local_default",
		Name:                  "Local Default Policy",
		Version:               "2026-03-17",
		RulesetDigest:         "sha256:local-default",
		DefaultVerdict:        policy.VerdictAllow,
		DefaultSandboxProfile: "local-default",
		BreakGlassAllowed:     true,
		Sandboxes: map[string]policy.SandboxProfile{
			"local-default": {Name: "local-default", Description: "safe local default"},
			"shell-safe":    {Name: "shell-safe", AllowShell: true, Description: "shell allowed in local process sandbox"},
			"patch-safe":    {Name: "patch-safe", AllowPatch: true, Description: "patch allowed in local filesystem sandbox"},
			"net-observe":   {Name: "net-observe", AllowExternalNet: true, Description: "network access allowed for governed boundary only"},
		},
		Capabilities: []policy.CapabilityPolicy{
			{Name: "tool.validate", Verdict: policy.VerdictAllow, SandboxProfile: "local-default", AllowedBoundaries: []policy.TrustBoundary{policy.TrustBoundaryLocalProcess}},
			{Name: "connector.invoke", Verdict: policy.VerdictAsk, SandboxProfile: "net-observe", AllowedBoundaries: []policy.TrustBoundary{policy.TrustBoundaryExternalNet}, RequireApproval: true},
			{Name: "mcp.ingest", Verdict: policy.VerdictAsk, SandboxProfile: "net-observe", AllowedBoundaries: []policy.TrustBoundary{policy.TrustBoundaryExternalNet}, RequireApproval: true},
		},
	}
}

func LoadBundle(path string) (policy.Bundle, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return policy.Bundle{}, fmt.Errorf("read policy bundle %s: %w", path, err)
	}
	var bundle policy.Bundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		return policy.Bundle{}, fmt.Errorf("decode policy bundle %s: %w", path, err)
	}
	if errs := policy.ValidateBundle(bundle); len(errs) > 0 {
		return policy.Bundle{}, fmt.Errorf("invalid policy bundle: %v", errs[0])
	}
	return bundle, nil
}

func (e *Engine) Evaluate(_ context.Context, req Request) (Result, error) {
	boundary := ResolveBoundary(req.Capability)
	rule, ok := matchCapability(e.bundle, req.Capability)
	profile, err := resolveSandbox(e.bundle, req.Capability)
	if err != nil {
		decision := newDecision(e.bundle, req, policy.VerdictDeny, boundary, "", err.Error(), e.now())
		_ = writeAudit(e.auditRoot, decision, e.now())
		return Result{Decision: decision}, nil
	}
	if !sandboxAllows(profile, req.Capability) {
		decision := newDecision(e.bundle, req, policy.VerdictDeny, boundary, profile.Name, "sandbox profile does not permit capability", e.now())
		_ = writeAudit(e.auditRoot, decision, e.now())
		return Result{Decision: decision}, nil
	}
	if ok && len(rule.AllowedBoundaries) > 0 && !containsBoundary(rule.AllowedBoundaries, boundary) {
		decision := newDecision(e.bundle, req, policy.VerdictDeny, boundary, profile.Name, "trust boundary not allowed by policy", e.now())
		_ = writeAudit(e.auditRoot, decision, e.now())
		return Result{Decision: decision}, nil
	}
	if verdict, reason := applyBreakGlass(e.bundle, rulePtr(rule, ok), req.BreakGlass); verdict != "" {
		decision := newDecision(e.bundle, req, verdict, boundary, profile.Name, reason, e.now())
		_ = writeAudit(e.auditRoot, decision, e.now())
		return Result{Decision: decision}, nil
	}
	verdict := permissionVerdict(e.bundle, req.Capability, req.RiskClass)
	reason := "policy default allow"
	if ok {
		reason = "capability policy matched"
	}
	switch verdict {
	case policy.VerdictAllow, policy.VerdictAsk, policy.VerdictDeny:
	default:
		verdict = policy.VerdictDeny
		reason = "unsupported policy verdict"
	}
	decision := newDecision(e.bundle, req, verdict, boundary, profile.Name, reason, e.now())
	if verdict == policy.VerdictAsk {
		decision.Reason = "operator checkpoint required - " + decision.Reason
	}
	if err := writeAudit(e.auditRoot, decision, e.now()); err != nil {
		return Result{}, err
	}
	return Result{Decision: decision}, nil
}

func (e *Engine) Bundle() policy.Bundle {
	return e.bundle
}

func matchCapability(bundle policy.Bundle, capability string) (policy.CapabilityPolicy, bool) {
	for _, rule := range bundle.Capabilities {
		if rule.Name == capability {
			return rule, true
		}
	}
	return policy.CapabilityPolicy{}, false
}

func containsBoundary(boundaries []policy.TrustBoundary, target policy.TrustBoundary) bool {
	for _, boundary := range boundaries {
		if boundary == target {
			return true
		}
	}
	return false
}

func rulePtr(rule policy.CapabilityPolicy, ok bool) *policy.CapabilityPolicy {
	if !ok {
		return nil
	}
	return &rule
}
