package abilities

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Rule 613.1d: Text-changing effects (Layer 3)
// This file implements effects that change the text of cards, including:
// - Word replacement (e.g., "Change all instances of 'target' to 'each'")
// - Color word changes (e.g., "Red becomes Blue")
// - Basic land type changes (applied in this layer)
// - Name changes

// ===== Text Change Types =====

// TextChangeType indicates what kind of text change is being applied
type TextChangeType int

const (
	TextChangeWord      TextChangeType = iota // Replace a word with another word
	TextChangeColor                           // Replace a color word with another color
	TextChangeBasicLand                       // Replace basic land type
	TextChangeName                            // Change the name of a card
	TextChangeSubtype                         // Replace creature/land subtype
	TextChangeAbility                         // Replace ability word
)

// ===== Text Changing Effect =====

// TextChangingEffect represents an effect that changes text on a card
// Rule 613.1d: Text-changing effects are applied in Layer 3
type TextChangingEffect struct {
	baseContinuousEffect
	sourceID    uuid.UUID
	changeType  TextChangeType
	fromText    string                      // Text to replace
	toText      string                      // Replacement text
	affectedIDs []uuid.UUID                 // Specific cards affected (empty = all)
	filter      func(card interface{}) bool // Filter for affected cards
	appliedTo   map[uuid.UUID]bool          // Track which cards this has been applied to
}

// NewTextChangingEffect creates a text-changing effect
func NewTextChangingEffect(
	source uuid.UUID,
	changeType TextChangeType,
	fromText, toText string,
	duration Duration,
) *TextChangingEffect {
	return &TextChangingEffect{
		baseContinuousEffect: baseContinuousEffect{
			layer:    LayerTextChanging,
			duration: duration,
		},
		sourceID:    source,
		changeType:  changeType,
		fromText:    fromText,
		toText:      toText,
		affectedIDs: make([]uuid.UUID, 0),
		appliedTo:   make(map[uuid.UUID]bool),
	}
}

// SetFilter sets a filter for which cards are affected
func (e *TextChangingEffect) SetFilter(filter func(card interface{}) bool) {
	e.filter = filter
}

// SetAffectedIDs sets specific cards to be affected
func (e *TextChangingEffect) SetAffectedIDs(ids []uuid.UUID) {
	e.affectedIDs = ids
}

// Applies checks if this effect applies to a given card
func (e *TextChangingEffect) Applies(cardID uuid.UUID, card interface{}) bool {
	// If specific IDs are set, only apply to those
	if len(e.affectedIDs) > 0 {
		for _, id := range e.affectedIDs {
			if id == cardID {
				return true
			}
		}
		return false
	}

	// If filter is set, use filter
	if e.filter != nil {
		return e.filter(card)
	}

	// Otherwise, apply to all
	return true
}

// Apply implements the Effect interface
func (e *TextChangingEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Text changes are applied during layer calculation
	// This method is called during continuous effect application
	return nil
}

// ApplyToCard applies the text change to a specific card
func (e *TextChangingEffect) ApplyToCard(card interface{}) error {
	// Implementation would modify the card's text based on changeType
	// This is called during Layer 3 calculation

	switch e.changeType {
	case TextChangeWord:
		return e.replaceWord(card)
	case TextChangeColor:
		return e.replaceColor(card)
	case TextChangeBasicLand:
		return e.replaceBasicLandType(card)
	case TextChangeName:
		return e.replaceName(card)
	case TextChangeSubtype:
		return e.replaceSubtype(card)
	case TextChangeAbility:
		return e.replaceAbilityWord(card)
	default:
		return fmt.Errorf("unknown text change type: %d", e.changeType)
	}
}

// replaceWord replaces occurrences of one word with another
func (e *TextChangingEffect) replaceWord(card interface{}) error {
	// Replace all instances of fromText with toText in rules text
	// This affects:
	// - Ability text
	// - Keyword abilities that reference words
	// Example: Magical Hack changes "plainswalk" to "islandwalk"
	return fmt.Errorf("word replacement not yet implemented")
}

// replaceColor replaces color words
func (e *TextChangingEffect) replaceColor(card interface{}) error {
	// Replace color words in text
	// Examples:
	// - "red" becomes "blue"
	// - "protection from red" becomes "protection from blue"
	return fmt.Errorf("color word replacement not yet implemented")
}

// replaceBasicLandType replaces basic land types
func (e *TextChangingEffect) replaceBasicLandType(card interface{}) error {
	// Replace basic land type in text
	// Examples:
	// - "Plains" becomes "Island"
	// - "Swampwalk" becomes "Islandwalk"
	return fmt.Errorf("basic land type replacement not yet implemented")
}

// replaceName changes a card's name
func (e *TextChangingEffect) replaceName(card interface{}) error {
	// Change the name of a card
	// Example: Volrath's Shapeshifter copies name
	return fmt.Errorf("name replacement not yet implemented")
}

// replaceSubtype replaces creature or land subtypes
func (e *TextChangingEffect) replaceSubtype(card interface{}) error {
	// Replace subtype references in text
	// Example: "Goblin" becomes "Elf" in all text
	return fmt.Errorf("subtype replacement not yet implemented")
}

// replaceAbilityWord replaces ability words
func (e *TextChangingEffect) replaceAbilityWord(card interface{}) error {
	// Replace ability words in text
	// Rare but possible with effects like Artificial Evolution
	return fmt.Errorf("ability word replacement not yet implemented")
}

// GetDescription returns a text description of the effect
func (e *TextChangingEffect) GetDescription() string {
	switch e.changeType {
	case TextChangeWord:
		return fmt.Sprintf("Change all instances of '%s' to '%s'", e.fromText, e.toText)
	case TextChangeColor:
		return fmt.Sprintf("Change all instances of the color word '%s' to '%s'", e.fromText, e.toText)
	case TextChangeBasicLand:
		return fmt.Sprintf("Change all instances of the basic land type '%s' to '%s'", e.fromText, e.toText)
	case TextChangeName:
		return fmt.Sprintf("Change name to '%s'", e.toText)
	case TextChangeSubtype:
		return fmt.Sprintf("Change all instances of '%s' to '%s'", e.fromText, e.toText)
	default:
		return "Change text"
	}
}

// ===== Specific Text-Changing Cards =====

// MagicalHackEffect represents effects like Magical Hack
// "{U} Instant - Change the text of target spell or permanent by replacing all instances
// of one basic land type with another. (For example, you may change 'swampwalk' to 'plainswalk')"
type MagicalHackEffect struct {
	*TextChangingEffect
	targetID uuid.UUID
}

// NewMagicalHackEffect creates a Magical Hack-style effect
func NewMagicalHackEffect(source, target uuid.UUID, fromLand, toLand string, duration Duration) *MagicalHackEffect {
	baseEffect := NewTextChangingEffect(source, TextChangeBasicLand, fromLand, toLand, duration)
	baseEffect.SetAffectedIDs([]uuid.UUID{target})

	return &MagicalHackEffect{
		TextChangingEffect: baseEffect,
		targetID:           target,
	}
}

// MindBendEffect represents effects like Mind Bend
// "{U} Instant - Change the text of target permanent by replacing all instances of one
// color word with another or one basic land type with another until end of turn."
type MindBendEffect struct {
	*TextChangingEffect
	targetID uuid.UUID
}

// NewMindBendEffect creates a Mind Bend-style effect
func NewMindBendEffect(source, target uuid.UUID, changeType TextChangeType, from, to string) *MindBendEffect {
	baseEffect := NewTextChangingEffect(source, changeType, from, to, DurationUntilEndOfTurn)
	baseEffect.SetAffectedIDs([]uuid.UUID{target})

	return &MindBendEffect{
		TextChangingEffect: baseEffect,
		targetID:           target,
	}
}

// ArtificialEvolutionEffect represents Artificial Evolution
// "{U} Instant - Change the text of target spell or permanent by replacing all instances
// of one creature type with another. The new creature type can't be a creature type that
// appears in that permanent's name."
type ArtificialEvolutionEffect struct {
	*TextChangingEffect
	targetID         uuid.UUID
	excludedSubtypes []string // Subtypes that can't be used (from card name)
}

// NewArtificialEvolutionEffect creates an Artificial Evolution-style effect
func NewArtificialEvolutionEffect(source, target uuid.UUID, fromSubtype, toSubtype string) *ArtificialEvolutionEffect {
	baseEffect := NewTextChangingEffect(source, TextChangeSubtype, fromSubtype, toSubtype, DurationPermanent)
	baseEffect.SetAffectedIDs([]uuid.UUID{target})

	return &ArtificialEvolutionEffect{
		TextChangingEffect: baseEffect,
		targetID:           target,
		excludedSubtypes:   make([]string, 0),
	}
}

// SetExcludedSubtypes sets which subtypes can't be used
func (e *ArtificialEvolutionEffect) SetExcludedSubtypes(subtypes []string) {
	e.excludedSubtypes = subtypes
}

// ===== Color Word Helper =====

// ColorWord represents the five Magic colors as text
type ColorWord string

const (
	ColorWordWhite ColorWord = "white"
	ColorWordBlue  ColorWord = "blue"
	ColorWordBlack ColorWord = "black"
	ColorWordRed   ColorWord = "red"
	ColorWordGreen ColorWord = "green"
)

// AllColorWords returns all color words
func AllColorWords() []ColorWord {
	return []ColorWord{
		ColorWordWhite,
		ColorWordBlue,
		ColorWordBlack,
		ColorWordRed,
		ColorWordGreen,
	}
}

// IsValidColorWord checks if a string is a valid color word
func IsValidColorWord(word string) bool {
	normalized := strings.ToLower(strings.TrimSpace(word))
	for _, cw := range AllColorWords() {
		if string(cw) == normalized {
			return true
		}
	}
	return false
}

// ===== Basic Land Type Helper =====

// BasicLandType represents the five basic land types
type BasicLandType string

const (
	BasicLandPlains   BasicLandType = "Plains"
	BasicLandIsland   BasicLandType = "Island"
	BasicLandSwamp    BasicLandType = "Swamp"
	BasicLandMountain BasicLandType = "Mountain"
	BasicLandForest   BasicLandType = "Forest"
)

// AllBasicLandTypes returns all basic land types
func AllBasicLandTypes() []BasicLandType {
	return []BasicLandType{
		BasicLandPlains,
		BasicLandIsland,
		BasicLandSwamp,
		BasicLandMountain,
		BasicLandForest,
	}
}

// IsValidBasicLandType checks if a string is a valid basic land type
func IsValidBasicLandType(landType string) bool {
	for _, blt := range AllBasicLandTypes() {
		if string(blt) == landType {
			return true
		}
	}
	return false
}

// GetLandwalkFromLandType returns the corresponding landwalk ability
func GetLandwalkFromLandType(landType BasicLandType) string {
	switch landType {
	case BasicLandPlains:
		return "plainswalk"
	case BasicLandIsland:
		return "islandwalk"
	case BasicLandSwamp:
		return "swampwalk"
	case BasicLandMountain:
		return "mountainwalk"
	case BasicLandForest:
		return "forestwalk"
	default:
		return ""
	}
}

// ===== Text Search and Replace Helper =====

// TextReplacementHelper provides utilities for text replacement
type TextReplacementHelper struct {
	caseSensitive bool
}

// NewTextReplacementHelper creates a helper for text replacement
func NewTextReplacementHelper(caseSensitive bool) *TextReplacementHelper {
	return &TextReplacementHelper{
		caseSensitive: caseSensitive,
	}
}

// ReplaceInText replaces all occurrences of one string with another
func (h *TextReplacementHelper) ReplaceInText(text, from, to string) string {
	if h.caseSensitive {
		return strings.ReplaceAll(text, from, to)
	}

	// Case-insensitive replacement while preserving original case
	// This is more complex - need to find all matches and replace individually
	lowerFrom := strings.ToLower(from)

	result := text
	index := 0
	for {
		pos := strings.Index(strings.ToLower(result[index:]), lowerFrom)
		if pos == -1 {
			break
		}
		actualPos := index + pos

		// Preserve capitalization of replacement
		original := result[actualPos : actualPos+len(from)]
		replacement := h.matchCase(original, to)

		result = result[:actualPos] + replacement + result[actualPos+len(from):]
		index = actualPos + len(replacement)
	}

	return result
}

// matchCase attempts to match the case pattern of the original
func (h *TextReplacementHelper) matchCase(original, replacement string) string {
	if len(original) == 0 || len(replacement) == 0 {
		return replacement
	}

	// If original is all uppercase, make replacement uppercase
	if original == strings.ToUpper(original) {
		return strings.ToUpper(replacement)
	}

	// If original is title case (first letter uppercase), make replacement title case
	if len(original) > 0 && original[0] == strings.ToUpper(string(original[0]))[0] {
		if len(replacement) == 0 {
			return replacement
		}
		return strings.ToUpper(string(replacement[0])) + strings.ToLower(replacement[1:])
	}

	// Otherwise, use lowercase
	return strings.ToLower(replacement)
}

// FindAllOccurrences finds all occurrences of a substring
func (h *TextReplacementHelper) FindAllOccurrences(text, search string) []int {
	positions := make([]int, 0)
	searchText := text
	if !h.caseSensitive {
		searchText = strings.ToLower(text)
		search = strings.ToLower(search)
	}

	index := 0
	for {
		pos := strings.Index(searchText[index:], search)
		if pos == -1 {
			break
		}
		positions = append(positions, index+pos)
		index += pos + 1
	}

	return positions
}

// ===== Integration with Existing Systems =====

// Text-changing effects integrate with:
// - Layer system: Applied in Layer 3 (after copy effects, control changes)
// - Continuous effects: Tracked as continuous effects
// - Card characteristics: Modify oracle text, not printed text
// - Targeting: Can target spells on the stack or permanents on battlefield

// Important rules:
// - Rule 613.1d: Text-changing effects are applied in Layer 3
// - Rule 612.1: Text-changing effects change only the text printed on the card
// - Rule 612.2: They don't affect card names unless explicitly stated
// - Rule 612.3: A continuous effect that modifies characteristics of a spell on the stack
//   continues to apply to the permanent that spell becomes

// Example usage for Magical Hack:
// 1. Cast Magical Hack targeting opponent's creature with Swampwalk
// 2. Choose to change "Swamp" to "Plains"
// 3. Create TextChangingEffect with:
//    - changeType: TextChangeBasicLand
//    - fromText: "Swamp"
//    - toText: "Plains"
//    - duration: DurationPermanent (changes last forever)
//    - affectedIDs: [target creature ID]
// 4. Effect is applied in Layer 3
// 5. "Swampwalk" becomes "Plainswalk"

// Example usage for Mind Bend:
// 1. Cast Mind Bend targeting permanent
// 2. Choose color word or land type change
// 3. Create TextChangingEffect with:
//    - duration: DurationUntilEndOfTurn
// 4. Effect expires at end of turn

// Example usage for Artificial Evolution:
// 1. Cast Artificial Evolution targeting permanent or spell
// 2. Choose creature type to change
// 3. Verify new type doesn't appear in card name
// 4. Create TextChangingEffect with:
//    - changeType: TextChangeSubtype
//    - duration: DurationPermanent
// 5. All instances of old subtype become new subtype

// ===== Common Text-Changing Card Examples =====

// Text-changing cards:
// - Magical Hack: Change one basic land type to another
// - Mind Bend: Change color word or basic land type until end of turn
// - Sleight of Mind: Change color word
// - Glamerdye: Change color word (with Conspire)
// - Trait Doctoring: Change color word or creature type
// - Artificial Evolution: Change creature type (permanent)
// - Prismatic Lace: Change all instances of one color word to another
// - Swirl the Mists: Global color word change each turn

// Effects that care about card text:
// - Cranial Insertion: Looks for creature types in text
// - Mistform Ultimus: Is all creature types (not text-changing, but related)
// - Shields of Velis Vel: Changes creature types globally

// Note: Text-changing effects are relatively rare in modern Magic but
// still appear occasionally, especially in supplemental sets
