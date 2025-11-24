package effects

import (
	"fmt"

	"github.com/google/uuid"
)

// Rule 707: Copy Effects
// Copy effects create a copy of an object and define what values are copied.
// Copies are created in Layer 1 of the layer system (Rule 613.1).

// CopyEffect represents an effect that makes one object a copy of another
// Examples: Clone, Copy Enchantment, Phantasmal Image, Progenitor Mimic
type CopyEffect struct {
	copyID        uuid.UUID          // The object becoming a copy
	originalID    uuid.UUID          // The object being copied
	modifications []CopyModification // Modifications to the copy
	duration      Duration
	description   string
	layer         int // Layer for continuous effect system (always 1 for copies)
}

// CopyModification represents a modification to a copy
// Example: "except it has flying", "except it's a 1/1"
type CopyModification interface {
	// Apply modifies the copied values
	Apply(copiedValues *CopiedValues) *CopiedValues

	// GetDescription returns a human-readable description
	GetDescription() string
}

// CopiedValues represents the copiable values of an object
// Rule 707.2: Copiable values are the values printed on the object,
// as modified by other copy effects
type CopiedValues struct {
	Name           string
	ManaCost       string
	Colors         []string
	ColorIndicator []string // For color indicator (Rule 204)
	Types          []string
	Subtypes       []string
	Supertypes     []string
	RulesText      string
	Power          int
	Toughness      int
	Loyalty        int
	Abilities      []interface{} // Abilities (would be []Ability in real implementation)

	// Special properties
	IsToken bool
	IsCopy  bool
	CopyOf  uuid.UUID // What this is a copy of

	// Modifications from other effects
	Modifications map[string]interface{}
}

// NewCopyEffect creates a copy effect
func NewCopyEffect(
	copyID, originalID uuid.UUID,
	modifications []CopyModification,
	duration Duration,
	description string,
) *CopyEffect {
	return &CopyEffect{
		copyID:        copyID,
		originalID:    originalID,
		modifications: modifications,
		duration:      duration,
		description:   description,
		layer:         1, // Copy effects are always Layer 1
	}
}

// GetCopyID returns the ID of the object becoming a copy
func (e *CopyEffect) GetCopyID() uuid.UUID {
	return e.copyID
}

// GetOriginalID returns the ID of the object being copied
func (e *CopyEffect) GetOriginalID() uuid.UUID {
	return e.originalID
}

// GetDuration returns the duration of this copy effect
func (e *CopyEffect) GetDuration() Duration {
	return e.duration
}

// GetLayer returns the layer (always 1 for copies)
func (e *CopyEffect) GetLayer() int {
	return e.layer
}

// Apply applies the copy effect to create copied values
func (e *CopyEffect) Apply(original *CopiedValues) *CopiedValues {
	// Start with the copiable values from the original
	copied := &CopiedValues{
		Name:           original.Name,
		ManaCost:       original.ManaCost,
		Colors:         append([]string{}, original.Colors...),
		ColorIndicator: append([]string{}, original.ColorIndicator...),
		Types:          append([]string{}, original.Types...),
		Subtypes:       append([]string{}, original.Subtypes...),
		Supertypes:     append([]string{}, original.Supertypes...),
		RulesText:      original.RulesText,
		Power:          original.Power,
		Toughness:      original.Toughness,
		Loyalty:        original.Loyalty,
		Abilities:      append([]interface{}{}, original.Abilities...),
		IsToken:        original.IsToken,
		IsCopy:         true,
		CopyOf:         e.originalID,
		Modifications:  make(map[string]interface{}),
	}

	// Apply modifications in order
	for _, mod := range e.modifications {
		copied = mod.Apply(copied)
	}

	return copied
}

// String returns a human-readable description
func (e *CopyEffect) String() string {
	if e.description != "" {
		return e.description
	}
	return fmt.Sprintf("Copy of %s", e.originalID.String())
}

// ===== Copy Modifications =====

// PowerToughnessModification changes the power/toughness of the copy
// Example: "except it's a 1/1"
type PowerToughnessModification struct {
	power       int
	toughness   int
	description string
}

// NewPowerToughnessModification creates a P/T modification
func NewPowerToughnessModification(power, toughness int, description string) *PowerToughnessModification {
	return &PowerToughnessModification{
		power:       power,
		toughness:   toughness,
		description: description,
	}
}

func (m *PowerToughnessModification) Apply(copied *CopiedValues) *CopiedValues {
	copied.Power = m.power
	copied.Toughness = m.toughness
	copied.Modifications["power_toughness"] = fmt.Sprintf("%d/%d", m.power, m.toughness)
	return copied
}

func (m *PowerToughnessModification) GetDescription() string {
	return m.description
}

// AddAbilityModification adds an ability to the copy
// Example: "except it has flying and haste"
type AddAbilityModification struct {
	abilities   []interface{} // Abilities to add (would be []Ability in real implementation)
	description string
}

// NewAddAbilityModification creates an ability addition modification
func NewAddAbilityModification(abilities []interface{}, description string) *AddAbilityModification {
	return &AddAbilityModification{
		abilities:   abilities,
		description: description,
	}
}

func (m *AddAbilityModification) Apply(copied *CopiedValues) *CopiedValues {
	copied.Abilities = append(copied.Abilities, m.abilities...)
	copied.Modifications["added_abilities"] = m.description
	return copied
}

func (m *AddAbilityModification) GetDescription() string {
	return m.description
}

// TypeModification changes the types of the copy
// Example: "in addition to its other types, it's a Zombie"
type TypeModification struct {
	addTypes       []string
	removeTypes    []string
	setTypes       []string // If set, replaces all types
	addSubtypes    []string
	removeSubtypes []string
	setSubtypes    []string // If set, replaces all subtypes
	description    string
}

// NewTypeModification creates a type modification
func NewTypeModification(
	addTypes, removeTypes, setTypes []string,
	addSubtypes, removeSubtypes, setSubtypes []string,
	description string,
) *TypeModification {
	return &TypeModification{
		addTypes:       addTypes,
		removeTypes:    removeTypes,
		setTypes:       setTypes,
		addSubtypes:    addSubtypes,
		removeSubtypes: removeSubtypes,
		setSubtypes:    setSubtypes,
		description:    description,
	}
}

func (m *TypeModification) Apply(copied *CopiedValues) *CopiedValues {
	// Handle types
	if len(m.setTypes) > 0 {
		copied.Types = append([]string{}, m.setTypes...)
	} else {
		// Remove types
		for _, removeType := range m.removeTypes {
			copied.Types = removeFromSlice(copied.Types, removeType)
		}
		// Add types
		for _, addType := range m.addTypes {
			if !contains(copied.Types, addType) {
				copied.Types = append(copied.Types, addType)
			}
		}
	}

	// Handle subtypes
	if len(m.setSubtypes) > 0 {
		copied.Subtypes = append([]string{}, m.setSubtypes...)
	} else {
		// Remove subtypes
		for _, removeSubtype := range m.removeSubtypes {
			copied.Subtypes = removeFromSlice(copied.Subtypes, removeSubtype)
		}
		// Add subtypes
		for _, addSubtype := range m.addSubtypes {
			if !contains(copied.Subtypes, addSubtype) {
				copied.Subtypes = append(copied.Subtypes, addSubtype)
			}
		}
	}

	copied.Modifications["type_change"] = m.description
	return copied
}

func (m *TypeModification) GetDescription() string {
	return m.description
}

// ColorModification changes the colors of the copy
// Example: "except it's red"
type ColorModification struct {
	colors       []string // New colors (replaces all colors)
	addColors    []string // Add these colors
	removeColors []string // Remove these colors
	description  string
}

// NewColorModification creates a color modification
func NewColorModification(colors, addColors, removeColors []string, description string) *ColorModification {
	return &ColorModification{
		colors:       colors,
		addColors:    addColors,
		removeColors: removeColors,
		description:  description,
	}
}

func (m *ColorModification) Apply(copied *CopiedValues) *CopiedValues {
	if len(m.colors) > 0 {
		// Set colors
		copied.Colors = append([]string{}, m.colors...)
	} else {
		// Remove colors
		for _, removeColor := range m.removeColors {
			copied.Colors = removeFromSlice(copied.Colors, removeColor)
		}
		// Add colors
		for _, addColor := range m.addColors {
			if !contains(copied.Colors, addColor) {
				copied.Colors = append(copied.Colors, addColor)
			}
		}
	}

	copied.Modifications["color_change"] = m.description
	return copied
}

func (m *ColorModification) GetDescription() string {
	return m.description
}

// NameModification changes the name of the copy
// Example: "except its name is [Name]"
type NameModification struct {
	name        string
	description string
}

// NewNameModification creates a name modification
func NewNameModification(name, description string) *NameModification {
	return &NameModification{
		name:        name,
		description: description,
	}
}

func (m *NameModification) Apply(copied *CopiedValues) *CopiedValues {
	copied.Name = m.name
	copied.Modifications["name_change"] = m.name
	return copied
}

func (m *NameModification) GetDescription() string {
	return m.description
}

// TokenCopyModification marks that the copy is a token
// Used when creating token copies (Rule 707.10)
type TokenCopyModification struct {
	description string
}

// NewTokenCopyModification creates a token copy modification
func NewTokenCopyModification(description string) *TokenCopyModification {
	return &TokenCopyModification{
		description: description,
	}
}

func (m *TokenCopyModification) Apply(copied *CopiedValues) *CopiedValues {
	copied.IsToken = true
	copied.Modifications["is_token"] = "true"
	return copied
}

func (m *TokenCopyModification) GetDescription() string {
	return m.description
}

// ===== Copy Manager =====

// CopyManager manages all copy effects in a game
// Handles Layer 1 of the layer system
type CopyManager struct {
	effects map[uuid.UUID]*CopyEffect   // copyID -> effect
	values  map[uuid.UUID]*CopiedValues // copyID -> current copied values
}

// NewCopyManager creates a new copy manager
func NewCopyManager() *CopyManager {
	return &CopyManager{
		effects: make(map[uuid.UUID]*CopyEffect),
		values:  make(map[uuid.UUID]*CopiedValues),
	}
}

// AddCopyEffect adds a copy effect
func (cm *CopyManager) AddCopyEffect(effect *CopyEffect) {
	cm.effects[effect.GetCopyID()] = effect
}

// RemoveCopyEffect removes a copy effect
func (cm *CopyManager) RemoveCopyEffect(copyID uuid.UUID) {
	delete(cm.effects, copyID)
	delete(cm.values, copyID)
}

// GetCopyEffect retrieves a copy effect
func (cm *CopyManager) GetCopyEffect(copyID uuid.UUID) (*CopyEffect, bool) {
	effect, ok := cm.effects[copyID]
	return effect, ok
}

// IsCopy checks if an object is currently a copy
func (cm *CopyManager) IsCopy(objectID uuid.UUID) bool {
	_, ok := cm.effects[objectID]
	return ok
}

// GetCopiedValues gets the current copied values for an object
func (cm *CopyManager) GetCopiedValues(copyID uuid.UUID) (*CopiedValues, bool) {
	values, ok := cm.values[copyID]
	return values, ok
}

// RecalculateCopy recalculates the copied values for an object
// This should be called when the original changes or in Layer 1
func (cm *CopyManager) RecalculateCopy(copyID uuid.UUID, originalValues *CopiedValues) {
	effect, ok := cm.effects[copyID]
	if !ok {
		return
	}

	// Apply the copy effect
	copiedValues := effect.Apply(originalValues)
	cm.values[copyID] = copiedValues
}

// RecalculateAll recalculates all copy effects
// Called during Layer 1 processing
func (cm *CopyManager) RecalculateAll(getOriginalValues func(uuid.UUID) *CopiedValues) {
	for copyID, effect := range cm.effects {
		originalValues := getOriginalValues(effect.GetOriginalID())
		if originalValues != nil {
			cm.RecalculateCopy(copyID, originalValues)
		}
	}
}

// Clear removes all copy effects
func (cm *CopyManager) Clear() {
	cm.effects = make(map[uuid.UUID]*CopyEffect)
	cm.values = make(map[uuid.UUID]*CopiedValues)
}

// ===== Helper Functions =====

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func removeFromSlice(slice []string, item string) []string {
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}

// ===== Spell/Ability Copy Support =====

// SpellCopy represents a copy of a spell on the stack
// Rule 707.10: Copies of spells can be created
type SpellCopy struct {
	copyID        uuid.UUID
	originalID    uuid.UUID
	controller    uuid.UUID
	targets       []uuid.UUID
	choices       map[string]interface{}
	modifications []CopyModification
}

// NewSpellCopy creates a copy of a spell
func NewSpellCopy(
	copyID, originalID, controller uuid.UUID,
	targets []uuid.UUID,
	choices map[string]interface{},
	modifications []CopyModification,
) *SpellCopy {
	return &SpellCopy{
		copyID:        copyID,
		originalID:    originalID,
		controller:    controller,
		targets:       targets,
		choices:       choices,
		modifications: modifications,
	}
}

// GetCopyID returns the copy's ID
func (sc *SpellCopy) GetCopyID() uuid.UUID {
	return sc.copyID
}

// GetOriginalID returns the original spell's ID
func (sc *SpellCopy) GetOriginalID() uuid.UUID {
	return sc.originalID
}

// GetController returns the controller of the copy
func (sc *SpellCopy) GetController() uuid.UUID {
	return sc.controller
}

// GetTargets returns the targets chosen for the copy
func (sc *SpellCopy) GetTargets() []uuid.UUID {
	return sc.targets
}

// GetChoices returns the choices made for the copy
func (sc *SpellCopy) GetChoices() map[string]interface{} {
	return sc.choices
}

// SpellCopyManager manages spell copies on the stack
type SpellCopyManager struct {
	copies map[uuid.UUID]*SpellCopy
}

// NewSpellCopyManager creates a new spell copy manager
func NewSpellCopyManager() *SpellCopyManager {
	return &SpellCopyManager{
		copies: make(map[uuid.UUID]*SpellCopy),
	}
}

// AddSpellCopy adds a spell copy
func (scm *SpellCopyManager) AddSpellCopy(copy *SpellCopy) {
	scm.copies[copy.GetCopyID()] = copy
}

// RemoveSpellCopy removes a spell copy (when it resolves or is countered)
func (scm *SpellCopyManager) RemoveSpellCopy(copyID uuid.UUID) {
	delete(scm.copies, copyID)
}

// GetSpellCopy retrieves a spell copy
func (scm *SpellCopyManager) GetSpellCopy(copyID uuid.UUID) (*SpellCopy, bool) {
	copy, ok := scm.copies[copyID]
	return copy, ok
}

// IsSpellCopy checks if a spell is a copy
func (scm *SpellCopyManager) IsSpellCopy(spellID uuid.UUID) bool {
	_, ok := scm.copies[spellID]
	return ok
}

// Clear removes all spell copies
func (scm *SpellCopyManager) Clear() {
	scm.copies = make(map[uuid.UUID]*SpellCopy)
}
