package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	ctools "srosv2/contracts/tools"
	"srosv2/internal/fabric/harness"
)

type Registry struct {
	root      string
	harness   *harness.Harness
	manifests map[string]ctools.Manifest
}

func New(root string, hr *harness.Harness, seed []ctools.Manifest) (*Registry, error) {
	if root == "" {
		root = filepath.Join("artifacts", "fabric", "registry")
	}
	r := &Registry{root: root, harness: hr, manifests: map[string]ctools.Manifest{}}
	files, err := listManifestFiles(root)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		manifest, err := loadManifest(file)
		if err != nil {
			return nil, err
		}
		r.manifests[manifest.Name] = manifest
	}
	for _, manifest := range seed {
		if _, exists := r.manifests[manifest.Name]; exists {
			continue
		}
		r.manifests[manifest.Name] = manifest
		if err := writeManifest(root, manifest); err != nil {
			return nil, err
		}
	}
	return r, nil
}

func (r *Registry) Root() string { return r.root }

func (r *Registry) List() []ctools.Manifest {
	names := make([]string, 0, len(r.manifests))
	for name := range r.manifests {
		names = append(names, name)
	}
	sort.Strings(names)
	items := make([]ctools.Manifest, 0, len(names))
	for _, name := range names {
		items = append(items, r.manifests[name])
	}
	return items
}

func (r *Registry) Get(name string) (ctools.Manifest, bool) {
	manifest, ok := r.manifests[strings.TrimSpace(name)]
	return manifest, ok
}

func (r *Registry) Validate(manifest ctools.Manifest) ctools.ValidationResult {
	return ValidateManifest(manifest, r.harness)
}

func (r *Registry) Register(manifest ctools.Manifest) (ctools.Manifest, ctools.ValidationResult, error) {
	result := r.Validate(manifest)
	if !result.Valid {
		return ctools.Manifest{}, result, fmt.Errorf("manifest validation failed")
	}
	manifest.Status = ctools.StateValidated
	r.manifests[manifest.Name] = manifest
	if err := writeManifest(r.root, manifest); err != nil {
		return ctools.Manifest{}, result, err
	}
	return manifest, result, nil
}

func (r *Registry) Admit(name string) (ctools.Manifest, error) {
	manifest, ok := r.Get(name)
	if !ok {
		return ctools.Manifest{}, fmt.Errorf("manifest %s not found", name)
	}
	if manifest.Status == ctools.StateDraft {
		manifest.Status = ctools.StateValidated
	}
	next := ctools.StateAdmitted
	if manifest.Experimental {
		next = ctools.StateExperimental
	}
	if manifest.Status == ctools.StateValidated {
		manifest.Status = next
	} else {
		var err error
		manifest, err = Transition(manifest, next)
		if err != nil {
			return ctools.Manifest{}, err
		}
	}
	r.manifests[manifest.Name] = manifest
	return manifest, writeManifest(r.root, manifest)
}

func (r *Registry) SetState(name string, to ctools.LifecycleState, reason string) (ctools.Manifest, error) {
	manifest, ok := r.Get(name)
	if !ok {
		return ctools.Manifest{}, fmt.Errorf("manifest %s not found", name)
	}
	if to == ctools.StateQuarantined {
		manifest.QuarantineReason = reason
	}
	updated, err := Transition(manifest, to)
	if err != nil {
		return ctools.Manifest{}, err
	}
	r.manifests[name] = updated
	return updated, writeManifest(r.root, updated)
}

func (r *Registry) Search(query ctools.SearchQuery) []ctools.SearchMatch {
	return Search(r.List(), query)
}

func (r *Registry) Shortlist(query ctools.SearchQuery) []ctools.SearchMatch {
	return Shortlist(r.List(), query)
}

func Reset(root string) error {
	return os.RemoveAll(root)
}
