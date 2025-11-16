package effects

// Duration represents how long an effect lasts
type Duration string

const (
	// DurationEndOfTurn - Effect expires at end of turn
	DurationEndOfTurn Duration = "EndOfTurn"
	
	// DurationEndOfCombat - Effect expires at end of combat
	DurationEndOfCombat Duration = "EndOfCombat"
	
	// DurationWhileOnBattlefield - Effect lasts while source is on battlefield
	DurationWhileOnBattlefield Duration = "WhileOnBattlefield"
	
	// DurationWhileControlled - Effect lasts while you control the source
	DurationWhileControlled Duration = "WhileControlled"
	
	// DurationUntilSourceLeaves - Effect lasts until source leaves battlefield
	DurationUntilSourceLeaves Duration = "UntilSourceLeaves"
	
	// DurationPermanent - Effect lasts indefinitely
	DurationPermanent Duration = "Permanent"
)

// EffectWithDuration represents an effect that has a duration
type EffectWithDuration interface {
	ContinuousEffect
	GetDuration() Duration
	GetSourceID() string
}

// CleanupEndOfCombatEffects removes effects that expire at end of combat
// Per Java: ContinuousEffects.removeEndOfCombatEffects()
func CleanupEndOfCombatEffects(system *LayerSystem) {
	if system == nil {
		return
	}
	
	system.mu.Lock()
	defer system.mu.Unlock()
	
	// Collect IDs of effects to remove
	var toRemove []string
	
	for layer, effectMap := range system.effects {
		for id, effect := range effectMap {
			// Check if effect has duration
			if durationEffect, ok := effect.(EffectWithDuration); ok {
				if durationEffect.GetDuration() == DurationEndOfCombat {
					toRemove = append(toRemove, id)
				}
			}
		}
		_ = layer // Suppress unused warning
	}
	
	// Remove expired effects
	for _, id := range toRemove {
		if layer, ok := system.index[id]; ok {
			delete(system.index, id)
			if layerMap, ok := system.effects[layer]; ok {
				delete(layerMap, id)
				if len(layerMap) == 0 {
					delete(system.effects, layer)
				}
			}
		}
	}
}

// CleanupEndOfTurnEffects removes effects that expire at end of turn
// Per Java: ContinuousEffects.removeEndOfTurnEffects()
func CleanupEndOfTurnEffects(system *LayerSystem) {
	if system == nil {
		return
	}
	
	system.mu.Lock()
	defer system.mu.Unlock()
	
	// Collect IDs of effects to remove
	var toRemove []string
	
	for layer, effectMap := range system.effects {
		for id, effect := range effectMap {
			// Check if effect has duration
			if durationEffect, ok := effect.(EffectWithDuration); ok {
				if durationEffect.GetDuration() == DurationEndOfTurn {
					toRemove = append(toRemove, id)
				}
			}
		}
		_ = layer // Suppress unused warning
	}
	
	// Remove expired effects
	for _, id := range toRemove {
		if layer, ok := system.index[id]; ok {
			delete(system.index, id)
			if layerMap, ok := system.effects[layer]; ok {
				delete(layerMap, id)
				if len(layerMap) == 0 {
					delete(system.effects, layer)
				}
			}
		}
	}
}

// CleanupSourceLeftBattlefieldEffects removes effects whose source left the battlefield
// Per Java: ContinuousEffects.removeInactiveEffects()
func CleanupSourceLeftBattlefieldEffects(system *LayerSystem, sourceID string) {
	if system == nil || sourceID == "" {
		return
	}
	
	system.mu.Lock()
	defer system.mu.Unlock()
	
	// Collect IDs of effects to remove
	var toRemove []string
	
	for layer, effectMap := range system.effects {
		for id, effect := range effectMap {
			// Check if effect depends on source being on battlefield
			if durationEffect, ok := effect.(EffectWithDuration); ok {
				if durationEffect.GetSourceID() == sourceID {
					duration := durationEffect.GetDuration()
					if duration == DurationWhileOnBattlefield || 
					   duration == DurationWhileControlled ||
					   duration == DurationUntilSourceLeaves {
						toRemove = append(toRemove, id)
					}
				}
			}
		}
		_ = layer // Suppress unused warning
	}
	
	// Remove expired effects
	for _, id := range toRemove {
		if layer, ok := system.index[id]; ok {
			delete(system.index, id)
			if layerMap, ok := system.effects[layer]; ok {
				delete(layerMap, id)
				if len(layerMap) == 0 {
					delete(system.effects, layer)
				}
			}
		}
	}
}
