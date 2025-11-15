package mana

// CostReduction represents a cost reduction effect.
type CostReduction struct {
	ID              string
	GenericReduction int
	ColoredReduction map[ManaType]int
	AppliesTo       func(cardID string, cost *ManaCost) bool // Function to check if reduction applies
}

// CostReductionManager manages cost reduction effects.
type CostReductionManager struct {
	reductions []*CostReduction
}

// NewCostReductionManager creates a new cost reduction manager.
func NewCostReductionManager() *CostReductionManager {
	return &CostReductionManager{
		reductions: make([]*CostReduction, 0),
	}
}

// AddReduction adds a cost reduction effect.
func (crm *CostReductionManager) AddReduction(reduction *CostReduction) {
	if reduction == nil {
		return
	}
	crm.reductions = append(crm.reductions, reduction)
}

// RemoveReduction removes a cost reduction effect by ID.
func (crm *CostReductionManager) RemoveReduction(id string) {
	for i, red := range crm.reductions {
		if red.ID == id {
			crm.reductions = append(crm.reductions[:i], crm.reductions[i+1:]...)
			return
		}
	}
}

// ApplyReductions applies all applicable cost reductions to a mana cost.
func (crm *CostReductionManager) ApplyReductions(cardID string, cost *ManaCost) *ManaCost {
	if cost == nil {
		return nil
	}

	reduced := cost
	totalGenericReduction := 0
	totalColoredReduction := make(map[ManaType]int)

	for _, reduction := range crm.reductions {
		if reduction.AppliesTo == nil || reduction.AppliesTo(cardID, cost) {
			totalGenericReduction += reduction.GenericReduction
			for mt, amount := range reduction.ColoredReduction {
				totalColoredReduction[mt] += amount
			}
		}
	}

	if totalGenericReduction > 0 || len(totalColoredReduction) > 0 {
		reduced = cost.ApplyReduction(totalGenericReduction, totalColoredReduction)
	}

	return reduced
}
