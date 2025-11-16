package effects

// EffectBuilder provides a fluent API for creating continuous effects
// This simplifies effect creation for card abilities
type EffectBuilder struct {
	sourceID  string
	targetIDs []string
	duration  Duration
}

// NewEffectBuilder creates a new effect builder
func NewEffectBuilder(sourceID string) *EffectBuilder {
	return &EffectBuilder{
		sourceID:  sourceID,
		targetIDs: []string{},
		duration:  DurationEndOfTurn,
	}
}

// Targeting sets the target IDs for the effect
func (b *EffectBuilder) Targeting(targetIDs ...string) *EffectBuilder {
	b.targetIDs = targetIDs
	return b
}

// UntilEndOfTurn sets the duration to end of turn
func (b *EffectBuilder) UntilEndOfTurn() *EffectBuilder {
	b.duration = DurationEndOfTurn
	return b
}

// UntilEndOfCombat sets the duration to end of combat
func (b *EffectBuilder) UntilEndOfCombat() *EffectBuilder {
	b.duration = DurationEndOfCombat
	return b
}

// WhileOnBattlefield sets the duration to while source is on battlefield
func (b *EffectBuilder) WhileOnBattlefield() *EffectBuilder {
	b.duration = DurationWhileOnBattlefield
	return b
}

// Permanent sets the duration to permanent
func (b *EffectBuilder) Permanent() *EffectBuilder {
	b.duration = DurationPermanent
	return b
}

// GrantAbility creates a GrantAbilityEffect
func (b *EffectBuilder) GrantAbility(abilityID string) *GrantAbilityEffect {
	return NewGrantAbilityEffect(b.sourceID, abilityID, b.targetIDs, b.duration)
}

// CantAttack creates a CantAttackEffect
func (b *EffectBuilder) CantAttack() *CantAttackEffect {
	return NewCantAttackEffect(b.sourceID, b.targetIDs, b.duration)
}

// CantBlock creates a CantBlockEffect
func (b *EffectBuilder) CantBlock() *CantBlockEffect {
	return NewCantBlockEffect(b.sourceID, b.targetIDs, b.duration)
}

// MustAttack creates a MustAttackEffect
func (b *EffectBuilder) MustAttack() *MustAttackEffect {
	return NewMustAttackEffect(b.sourceID, b.targetIDs, b.duration)
}

// EffectManager provides high-level effect management operations
type EffectManager struct {
	layerSystem *LayerSystem
}

// NewEffectManager creates a new effect manager
func NewEffectManager(layerSystem *LayerSystem) *EffectManager {
	return &EffectManager{
		layerSystem: layerSystem,
	}
}

// AddEffect adds an effect to the layer system
func (m *EffectManager) AddEffect(effect ContinuousEffect) string {
	if m.layerSystem == nil {
		return ""
	}
	return m.layerSystem.AddEffect(effect)
}

// RemoveEffect removes an effect by ID
func (m *EffectManager) RemoveEffect(effectID string) {
	if m.layerSystem != nil {
		m.layerSystem.RemoveEffect(effectID)
	}
}

// RemoveEffectsFromSource removes all effects from a specific source
func (m *EffectManager) RemoveEffectsFromSource(sourceID string) {
	if m.layerSystem == nil || sourceID == "" {
		return
	}
	
	// Collect effect IDs to remove
	var toRemove []string
	
	m.layerSystem.mu.RLock()
	for _, layerMap := range m.layerSystem.effects {
		for id, effect := range layerMap {
			if durationEffect, ok := effect.(EffectWithDuration); ok {
				if durationEffect.GetSourceID() == sourceID {
					toRemove = append(toRemove, id)
				}
			}
		}
	}
	m.layerSystem.mu.RUnlock()
	
	// Remove effects
	for _, id := range toRemove {
		m.layerSystem.RemoveEffect(id)
	}
}

// GetEffectsForCard returns all effects affecting a card
func (m *EffectManager) GetEffectsForCard(cardID string) []ContinuousEffect {
	if m.layerSystem == nil {
		return nil
	}
	return m.layerSystem.GetEffectsForCard(cardID)
}

// HasCantAttackEffect checks if a card has a can't attack effect
func (m *EffectManager) HasCantAttackEffect(cardID string) bool {
	if m.layerSystem == nil {
		return false
	}
	return m.layerSystem.HasEffectType(cardID, func(effect ContinuousEffect) bool {
		_, ok := effect.(*CantAttackEffect)
		return ok
	})
}

// HasCantBlockEffect checks if a card has a can't block effect
func (m *EffectManager) HasCantBlockEffect(cardID string) bool {
	if m.layerSystem == nil {
		return false
	}
	return m.layerSystem.HasEffectType(cardID, func(effect ContinuousEffect) bool {
		_, ok := effect.(*CantBlockEffect)
		return ok
	})
}

// HasMustAttackEffect checks if a card has a must attack effect
func (m *EffectManager) HasMustAttackEffect(cardID string) bool {
	if m.layerSystem == nil {
		return false
	}
	return m.layerSystem.HasEffectType(cardID, func(effect ContinuousEffect) bool {
		_, ok := effect.(*MustAttackEffect)
		return ok
	})
}

// HasGrantedAbility checks if a card has a specific granted ability
func (m *EffectManager) HasGrantedAbility(cardID, abilityID string) bool {
	if m.layerSystem == nil {
		return false
	}
	return m.layerSystem.HasEffectType(cardID, func(effect ContinuousEffect) bool {
		if grantEffect, ok := effect.(*GrantAbilityEffect); ok {
			return grantEffect.GetAbilityID() == abilityID
		}
		return false
	})
}
