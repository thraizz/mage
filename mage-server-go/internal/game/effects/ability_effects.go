package effects

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// GrantAbilityEffect grants an ability to target permanents
// Per Java GainAbilityTargetEffect
type GrantAbilityEffect struct {
	id           string
	sourceID     string
	abilityID    string // The ability being granted (e.g., "FlyingAbility")
	targetIDs    []string // IDs of affected permanents
	duration     Duration // How long the effect lasts
}

// NewGrantAbilityEffect creates a new ability-granting effect
func NewGrantAbilityEffect(sourceID, abilityID string, targetIDs []string, duration Duration) *GrantAbilityEffect {
	source := strings.TrimSpace(sourceID)
	ability := strings.TrimSpace(abilityID)
	seed := fmt.Sprintf("%s|%s|%v|%s", source, ability, targetIDs, duration)
	id := uuid.NewSHA1(uuid.NameSpaceOID, []byte(seed)).String()

	return &GrantAbilityEffect{
		id:        id,
		sourceID:  source,
		abilityID: ability,
		targetIDs: targetIDs,
		duration:  duration,
	}
}

// GetDuration returns the duration of the effect
func (e *GrantAbilityEffect) GetDuration() Duration {
	return e.duration
}

// GetSourceID returns the source ID of the effect
func (e *GrantAbilityEffect) GetSourceID() string {
	return e.sourceID
}

// ID returns the unique identifier
func (e *GrantAbilityEffect) ID() string {
	return e.id
}

// Layer identifies the layer in which the effect applies (Layer 6 - Ability)
func (e *GrantAbilityEffect) Layer() Layer {
	return LayerAbility
}

// AppliesTo determines whether the snapshot should receive the ability
func (e *GrantAbilityEffect) AppliesTo(snapshot *Snapshot) bool {
	if snapshot == nil {
		return false
	}
	
	// Check if this snapshot is one of the targets
	for _, targetID := range e.targetIDs {
		if snapshot.CardID == targetID {
			return true
		}
	}
	
	return false
}

// Apply mutates the snapshot (note: actual ability granting would need more infrastructure)
// For now, this is a placeholder that demonstrates the pattern
func (e *GrantAbilityEffect) Apply(snapshot *Snapshot) {
	// TODO: Implement ability granting
	// This would need to interact with the card's ability list
	// For now, this is a structural placeholder
}

// GetAbilityID returns the ability being granted
func (e *GrantAbilityEffect) GetAbilityID() string {
	return e.abilityID
}

// GetTargetIDs returns the target IDs
func (e *GrantAbilityEffect) GetTargetIDs() []string {
	return e.targetIDs
}

// CantAttackEffect prevents a creature from attacking
// Per Java CantAttackTargetEffect / RestrictionEffect
type CantAttackEffect struct {
	id        string
	sourceID  string
	targetIDs []string
	duration  Duration
}

// NewCantAttackEffect creates a new can't attack effect
func NewCantAttackEffect(sourceID string, targetIDs []string, duration Duration) *CantAttackEffect {
	source := strings.TrimSpace(sourceID)
	seed := fmt.Sprintf("%s|cant-attack|%v|%s", source, targetIDs, duration)
	id := uuid.NewSHA1(uuid.NameSpaceOID, []byte(seed)).String()

	return &CantAttackEffect{
		id:        id,
		sourceID:  source,
		targetIDs: targetIDs,
		duration:  duration,
	}
}

// GetDuration returns the duration of the effect
func (e *CantAttackEffect) GetDuration() Duration {
	return e.duration
}

// GetSourceID returns the source ID of the effect
func (e *CantAttackEffect) GetSourceID() string {
	return e.sourceID
}

// ID returns the unique identifier
func (e *CantAttackEffect) ID() string {
	return e.id
}

// Layer identifies the layer (Layer 9 - Rules Effects)
// Note: This is a conceptual layer - actual implementation would be in combat system
func (e *CantAttackEffect) Layer() Layer {
	return LayerPowerToughness // Using PT layer as placeholder since we only have 7 layers
}

// AppliesTo determines whether the snapshot is affected
func (e *CantAttackEffect) AppliesTo(snapshot *Snapshot) bool {
	if snapshot == nil {
		return false
	}
	
	for _, targetID := range e.targetIDs {
		if snapshot.CardID == targetID {
			return true
		}
	}
	
	return false
}

// Apply is a no-op for restriction effects (they're checked during combat)
func (e *CantAttackEffect) Apply(snapshot *Snapshot) {
	// Restriction effects don't modify the snapshot
	// They're checked by the combat system when declaring attackers
}

// GetTargetIDs returns the target IDs
func (e *CantAttackEffect) GetTargetIDs() []string {
	return e.targetIDs
}

// CantBlockEffect prevents a creature from blocking
// Per Java CantBlockTargetEffect / RestrictionEffect
type CantBlockEffect struct {
	id        string
	sourceID  string
	targetIDs []string
	duration  Duration
}

// NewCantBlockEffect creates a new can't block effect
func NewCantBlockEffect(sourceID string, targetIDs []string, duration Duration) *CantBlockEffect {
	source := strings.TrimSpace(sourceID)
	seed := fmt.Sprintf("%s|cant-block|%v|%s", source, targetIDs, duration)
	id := uuid.NewSHA1(uuid.NameSpaceOID, []byte(seed)).String()

	return &CantBlockEffect{
		id:        id,
		sourceID:  source,
		targetIDs: targetIDs,
		duration:  duration,
	}
}

// GetDuration returns the duration of the effect
func (e *CantBlockEffect) GetDuration() Duration {
	return e.duration
}

// GetSourceID returns the source ID of the effect
func (e *CantBlockEffect) GetSourceID() string {
	return e.sourceID
}

// ID returns the unique identifier
func (e *CantBlockEffect) ID() string {
	return e.id
}

// Layer identifies the layer (Rules Effects)
func (e *CantBlockEffect) Layer() Layer {
	return LayerPowerToughness // Using PT layer as placeholder
}

// AppliesTo determines whether the snapshot is affected
func (e *CantBlockEffect) AppliesTo(snapshot *Snapshot) bool {
	if snapshot == nil {
		return false
	}
	
	for _, targetID := range e.targetIDs {
		if snapshot.CardID == targetID {
			return true
		}
	}
	
	return false
}

// Apply is a no-op for restriction effects
func (e *CantBlockEffect) Apply(snapshot *Snapshot) {
	// Restriction effects don't modify the snapshot
	// They're checked by the combat system when declaring blockers
}

// GetTargetIDs returns the target IDs
func (e *CantBlockEffect) GetTargetIDs() []string {
	return e.targetIDs
}

// MustAttackEffect requires a creature to attack if able
// Per Java AttacksIfAbleSourceEffect / RequirementEffect
type MustAttackEffect struct {
	id        string
	sourceID  string
	targetIDs []string
	duration  Duration
}

// NewMustAttackEffect creates a new must attack effect
func NewMustAttackEffect(sourceID string, targetIDs []string, duration Duration) *MustAttackEffect {
	source := strings.TrimSpace(sourceID)
	seed := fmt.Sprintf("%s|must-attack|%v|%s", source, targetIDs, duration)
	id := uuid.NewSHA1(uuid.NameSpaceOID, []byte(seed)).String()

	return &MustAttackEffect{
		id:        id,
		sourceID:  source,
		targetIDs: targetIDs,
		duration:  duration,
	}
}

// GetDuration returns the duration of the effect
func (e *MustAttackEffect) GetDuration() Duration {
	return e.duration
}

// GetSourceID returns the source ID of the effect
func (e *MustAttackEffect) GetSourceID() string {
	return e.sourceID
}

// ID returns the unique identifier
func (e *MustAttackEffect) ID() string {
	return e.id
}

// Layer identifies the layer (Rules Effects)
func (e *MustAttackEffect) Layer() Layer {
	return LayerPowerToughness // Using PT layer as placeholder
}

// AppliesTo determines whether the snapshot is affected
func (e *MustAttackEffect) AppliesTo(snapshot *Snapshot) bool {
	if snapshot == nil {
		return false
	}
	
	for _, targetID := range e.targetIDs {
		if snapshot.CardID == targetID {
			return true
		}
	}
	
	return false
}

// Apply is a no-op for requirement effects
func (e *MustAttackEffect) Apply(snapshot *Snapshot) {
	// Requirement effects don't modify the snapshot
	// They're checked by the combat system when declaring attackers
}

// GetTargetIDs returns the target IDs
func (e *MustAttackEffect) GetTargetIDs() []string {
	return e.targetIDs
}
