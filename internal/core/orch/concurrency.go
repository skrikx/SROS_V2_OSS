package orch

import "srosv2/contracts/runcontract"

type ConcurrencyRules struct {
	MaxParallel int  `json:"max_parallel"`
	AllowSpawn  bool `json:"allow_spawn"`
}

func RulesForRisk(risk runcontract.RiskClass) ConcurrencyRules {
	switch risk {
	case runcontract.RiskClassCritical:
		return ConcurrencyRules{MaxParallel: 1, AllowSpawn: false}
	case runcontract.RiskClassHigh:
		return ConcurrencyRules{MaxParallel: 1, AllowSpawn: false}
	case runcontract.RiskClassMedium:
		return ConcurrencyRules{MaxParallel: 2, AllowSpawn: true}
	default:
		return ConcurrencyRules{MaxParallel: 3, AllowSpawn: true}
	}
}
