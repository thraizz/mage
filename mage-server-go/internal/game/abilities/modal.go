package abilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Mode represents a single mode of a modal spell or ability
// MTG Rule 700.2: Modal spells and abilities require the player to choose one or more modes
type Mode struct {
	ID          int                // Numeric ID for this mode (1, 2, 3, etc.)
	Text        string             // Mode text (e.g., "Counter target spell")
	Effects     []Effect           // Effects for this mode
	Targets     *TargetRequirement // Target requirements for this mode
	Restriction ModeRestriction    // Conditions for selecting this mode
}

// ModeRestriction defines when a mode can be chosen
type ModeRestriction interface {
	// CanChoose checks if this mode can be chosen
	CanChoose(ctx context.Context, game GameContext, source uuid.UUID) bool

	// GetDescription returns a human-readable description
	GetDescription() string
}

// AlwaysAllowedRestriction allows the mode to always be chosen
type AlwaysAllowedRestriction struct{}

func (r *AlwaysAllowedRestriction) CanChoose(ctx context.Context, game GameContext, source uuid.UUID) bool {
	return true
}

func (r *AlwaysAllowedRestriction) GetDescription() string {
	return ""
}

// ConditionalRestriction restricts mode selection based on game state
type ConditionalRestriction struct {
	Condition func(ctx context.Context, game GameContext, source uuid.UUID) bool
	Text      string
}

func (r *ConditionalRestriction) CanChoose(ctx context.Context, game GameContext, source uuid.UUID) bool {
	return r.Condition(ctx, game, source)
}

func (r *ConditionalRestriction) GetDescription() string {
	return r.Text
}

// ModalSpellAbility represents a spell with multiple modes
// Examples: Cryptic Command, Charm cycle, Command cycle
type ModalSpellAbility struct {
	baseAbility
	ManaCost    *ManaCost
	Modes       []*Mode
	MinModes    int  // Minimum number of modes to choose (usually 1)
	MaxModes    int  // Maximum number of modes to choose (usually 1)
	Unique      bool // If true, can't choose the same mode twice
	ChosenModes []int
}

// NewModalSpellAbility creates a new modal spell ability
func NewModalSpellAbility(sourceID uuid.UUID, manaCost string, modes []*Mode, minModes, maxModes int) (*ModalSpellAbility, error) {
	cost, err := ParseManaCost(manaCost)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mana cost: %w", err)
	}

	// Build text from modes
	text := "Choose "
	if minModes == maxModes {
		if minModes == 1 {
			text += "one —"
		} else {
			text += fmt.Sprintf("%d —", minModes)
		}
	} else {
		text += fmt.Sprintf("up to %d —", maxModes)
	}

	return &ModalSpellAbility{
		baseAbility: newBaseAbility(sourceID, text),
		ManaCost:    cost,
		Modes:       modes,
		MinModes:    minModes,
		MaxModes:    maxModes,
		Unique:      true, // Default to unique mode selection
	}, nil
}

// GetType returns the ability type
func (a *ModalSpellAbility) GetType() AbilityType {
	return AbilityTypeSpell
}

// CanActivate checks if this modal spell can be cast
func (a *ModalSpellAbility) CanActivate(ctx context.Context, game GameContext) bool {
	// Check if at least MinModes are available to choose
	availableModes := 0
	for _, mode := range a.Modes {
		if mode.Restriction == nil || mode.Restriction.CanChoose(ctx, game, a.sourceID) {
			availableModes++
		}
	}

	return availableModes >= a.MinModes
}

// SetChosenModes sets which modes were chosen when casting
func (a *ModalSpellAbility) SetChosenModes(modeIDs []int) error {
	// Validate mode count
	if len(modeIDs) < a.MinModes {
		return fmt.Errorf("must choose at least %d modes", a.MinModes)
	}
	if len(modeIDs) > a.MaxModes {
		return fmt.Errorf("can choose at most %d modes", a.MaxModes)
	}

	// Validate uniqueness if required
	if a.Unique {
		seen := make(map[int]bool)
		for _, id := range modeIDs {
			if seen[id] {
				return fmt.Errorf("cannot choose the same mode twice")
			}
			seen[id] = true
		}
	}

	// Validate mode IDs exist
	for _, id := range modeIDs {
		found := false
		for _, mode := range a.Modes {
			if mode.ID == id {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid mode ID: %d", id)
		}
	}

	a.ChosenModes = modeIDs
	return nil
}

// GetChosenModes returns the IDs of chosen modes
func (a *ModalSpellAbility) GetChosenModes() []int {
	return a.ChosenModes
}

// GetMode returns a mode by its ID
func (a *ModalSpellAbility) GetMode(id int) (*Mode, error) {
	for _, mode := range a.Modes {
		if mode.ID == id {
			return mode, nil
		}
	}
	return nil, fmt.Errorf("mode %d not found", id)
}

// GetAllModes returns all modes
func (a *ModalSpellAbility) GetAllModes() []*Mode {
	return a.Modes
}

// GetAvailableModes returns modes that can currently be chosen
func (a *ModalSpellAbility) GetAvailableModes(ctx context.Context, game GameContext) []*Mode {
	available := make([]*Mode, 0, len(a.Modes))
	for _, mode := range a.Modes {
		if mode.Restriction == nil || mode.Restriction.CanChoose(ctx, game, a.sourceID) {
			available = append(available, mode)
		}
	}
	return available
}

// Resolve resolves this modal spell by applying effects from chosen modes
func (a *ModalSpellAbility) Resolve(ctx context.Context, game GameContext) error {
	if len(a.ChosenModes) == 0 {
		return fmt.Errorf("no modes chosen")
	}

	// Apply effects for each chosen mode in order
	for _, modeID := range a.ChosenModes {
		mode, err := a.GetMode(modeID)
		if err != nil {
			return fmt.Errorf("failed to get mode %d: %w", modeID, err)
		}

		// TODO: Get targets for this specific mode
		targets := []uuid.UUID{}

		// Apply each effect in this mode
		for _, effect := range mode.Effects {
			if err := effect.Apply(ctx, game, a.sourceID, targets); err != nil {
				return fmt.Errorf("failed to apply effect for mode %d: %w", modeID, err)
			}
		}
	}

	return nil
}

// GetManaCost returns the mana cost
func (a *ModalSpellAbility) GetManaCost() *ManaCost {
	return a.ManaCost
}

// GetTargetsForMode returns target requirements for a specific mode
func (a *ModalSpellAbility) GetTargetsForMode(modeID int) (*TargetRequirement, error) {
	mode, err := a.GetMode(modeID)
	if err != nil {
		return nil, err
	}
	return mode.Targets, nil
}

// GetAllTargetRequirements returns target requirements for all chosen modes
func (a *ModalSpellAbility) GetAllTargetRequirements() []*TargetRequirement {
	requirements := make([]*TargetRequirement, 0, len(a.ChosenModes))
	for _, modeID := range a.ChosenModes {
		if mode, err := a.GetMode(modeID); err == nil && mode.Targets != nil {
			requirements = append(requirements, mode.Targets)
		}
	}
	return requirements
}

// String returns a human-readable description
func (a *ModalSpellAbility) String() string {
	text := a.text + "\n"
	for i, mode := range a.Modes {
		bullet := "•"
		text += fmt.Sprintf("%s %s\n", bullet, mode.Text)
		if i < len(a.Modes)-1 {
			text += "\n"
		}
	}
	return text
}

// ModalSpellBuilder helps build modal spells fluently
type ModalSpellBuilder struct {
	sourceID uuid.UUID
	manaCost string
	modes    []*Mode
	minModes int
	maxModes int
	unique   bool
}

// NewModalSpellBuilder creates a new builder
func NewModalSpellBuilder(sourceID uuid.UUID, manaCost string) *ModalSpellBuilder {
	return &ModalSpellBuilder{
		sourceID: sourceID,
		manaCost: manaCost,
		modes:    make([]*Mode, 0),
		minModes: 1,
		maxModes: 1,
		unique:   true,
	}
}

// SetModeRange sets the min/max number of modes
func (b *ModalSpellBuilder) SetModeRange(min, max int) *ModalSpellBuilder {
	b.minModes = min
	b.maxModes = max
	return b
}

// AllowDuplicateModes allows choosing the same mode multiple times
func (b *ModalSpellBuilder) AllowDuplicateModes() *ModalSpellBuilder {
	b.unique = false
	return b
}

// AddMode adds a mode to the spell
func (b *ModalSpellBuilder) AddMode(text string, effects []Effect) *ModalSpellBuilder {
	mode := &Mode{
		ID:          len(b.modes) + 1,
		Text:        text,
		Effects:     effects,
		Restriction: &AlwaysAllowedRestriction{},
	}
	b.modes = append(b.modes, mode)
	return b
}

// AddModeWithTargets adds a mode with target requirements
func (b *ModalSpellBuilder) AddModeWithTargets(text string, effects []Effect, targets *TargetRequirement) *ModalSpellBuilder {
	mode := &Mode{
		ID:          len(b.modes) + 1,
		Text:        text,
		Effects:     effects,
		Targets:     targets,
		Restriction: &AlwaysAllowedRestriction{},
	}
	b.modes = append(b.modes, mode)
	return b
}

// AddModeWithRestriction adds a mode with a restriction
func (b *ModalSpellBuilder) AddModeWithRestriction(text string, effects []Effect, restriction ModeRestriction) *ModalSpellBuilder {
	mode := &Mode{
		ID:          len(b.modes) + 1,
		Text:        text,
		Effects:     effects,
		Restriction: restriction,
	}
	b.modes = append(b.modes, mode)
	return b
}

// Build creates the modal spell ability
func (b *ModalSpellBuilder) Build() (*ModalSpellAbility, error) {
	if len(b.modes) == 0 {
		return nil, fmt.Errorf("modal spell must have at least one mode")
	}

	if b.minModes > len(b.modes) {
		return nil, fmt.Errorf("minModes (%d) cannot exceed number of modes (%d)", b.minModes, len(b.modes))
	}

	if b.maxModes > len(b.modes) {
		return nil, fmt.Errorf("maxModes (%d) cannot exceed number of modes (%d)", b.maxModes, len(b.modes))
	}

	if b.minModes > b.maxModes {
		return nil, fmt.Errorf("minModes (%d) cannot exceed maxModes (%d)", b.minModes, b.maxModes)
	}

	ability, err := NewModalSpellAbility(b.sourceID, b.manaCost, b.modes, b.minModes, b.maxModes)
	if err != nil {
		return nil, err
	}

	ability.Unique = b.unique
	return ability, nil
}
