package release

type Stage string

const (
	StageDraft      Stage = "draft"
	StageValidated  Stage = "validated"
	StagePromoted   Stage = "promoted"
	StageRolledBack Stage = "rolled_back"
)
