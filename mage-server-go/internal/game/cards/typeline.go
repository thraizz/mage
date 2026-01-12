package cards

import (
	"strings"
)

// ParsedTypeLine contains the parsed components of a card's type line
type ParsedTypeLine struct {
	Supertypes []string // e.g., "Legendary", "Basic", "Snow"
	Types      []string // e.g., "Creature", "Instant", "Land"
	Subtypes   []string // e.g., "Human", "Wizard", "Plains"
}

// Known supertypes in Magic: The Gathering
var knownSupertypes = map[string]bool{
	"LEGENDARY": true,
	"BASIC":     true,
	"SNOW":      true,
	"WORLD":     true,
	"ONGOING":   true,
	"ELITE":     true,
	"HOST":      true,
}

// Known card types in Magic: The Gathering
var knownTypes = map[string]bool{
	"ARTIFACT":     true,
	"CREATURE":     true,
	"ENCHANTMENT":  true,
	"INSTANT":      true,
	"LAND":         true,
	"PLANESWALKER": true,
	"SORCERY":      true,
	"TRIBAL":       true,
	"CONSPIRACY":   true,
	"PHENOMENON":   true,
	"PLANE":        true,
	"SCHEME":       true,
	"VANGUARD":     true,
	"BATTLE":       true,
	"DUNGEON":      true,
	"KINDRED":      true,
}

// ParseTypeLine parses a Magic card type line into its components
// Examples:
//   - "Legendary Creature — Human Wizard" → Legendary | Creature | [Human, Wizard]
//   - "Instant" → [] | Instant | []
//   - "Basic Land — Forest" → Basic | Land | [Forest]
//   - "Artifact Creature — Golem" → [] | [Artifact, Creature] | [Golem]
func ParseTypeLine(typeLine string) ParsedTypeLine {
	result := ParsedTypeLine{
		Supertypes: []string{},
		Types:      []string{},
		Subtypes:   []string{},
	}

	if typeLine == "" {
		return result
	}

	// Split on em dash (—) or regular dash to separate types from subtypes
	// Some cards use different dash characters
	var mainPart, subtypePart string
	if idx := strings.Index(typeLine, "—"); idx >= 0 {
		mainPart = strings.TrimSpace(typeLine[:idx])
		subtypePart = strings.TrimSpace(typeLine[idx+len("—"):])
	} else if idx := strings.Index(typeLine, " - "); idx >= 0 {
		mainPart = strings.TrimSpace(typeLine[:idx])
		subtypePart = strings.TrimSpace(typeLine[idx+3:])
	} else if idx := strings.Index(typeLine, " — "); idx >= 0 {
		mainPart = strings.TrimSpace(typeLine[:idx])
		subtypePart = strings.TrimSpace(typeLine[idx+3:])
	} else {
		mainPart = typeLine
	}

	// Parse main part (supertypes and types)
	words := strings.Fields(mainPart)
	for _, word := range words {
		upperWord := strings.ToUpper(word)

		if knownSupertypes[upperWord] {
			result.Supertypes = append(result.Supertypes, upperWord)
		} else if knownTypes[upperWord] {
			result.Types = append(result.Types, upperWord)
		} else {
			// Unknown type - add as regular type
			result.Types = append(result.Types, upperWord)
		}
	}

	// Parse subtypes
	if subtypePart != "" {
		subtypes := strings.Fields(subtypePart)
		for _, subtype := range subtypes {
			result.Subtypes = append(result.Subtypes, strings.ToUpper(subtype))
		}
	}

	return result
}

// IsCreature returns true if the type line represents a creature
func (p ParsedTypeLine) IsCreature() bool {
	return contains(p.Types, "CREATURE")
}

// IsInstant returns true if the type line represents an instant
func (p ParsedTypeLine) IsInstant() bool {
	return contains(p.Types, "INSTANT")
}

// IsSorcery returns true if the type line represents a sorcery
func (p ParsedTypeLine) IsSorcery() bool {
	return contains(p.Types, "SORCERY")
}

// IsEnchantment returns true if the type line represents an enchantment
func (p ParsedTypeLine) IsEnchantment() bool {
	return contains(p.Types, "ENCHANTMENT")
}

// IsArtifact returns true if the type line represents an artifact
func (p ParsedTypeLine) IsArtifact() bool {
	return contains(p.Types, "ARTIFACT")
}

// IsPlaneswalker returns true if the type line represents a planeswalker
func (p ParsedTypeLine) IsPlaneswalker() bool {
	return contains(p.Types, "PLANESWALKER")
}

// IsLand returns true if the type line represents a land
func (p ParsedTypeLine) IsLand() bool {
	return contains(p.Types, "LAND")
}

// IsLegendary returns true if the type line includes the Legendary supertype
func (p ParsedTypeLine) IsLegendary() bool {
	return contains(p.Supertypes, "LEGENDARY")
}

// IsBasic returns true if the type line includes the Basic supertype
func (p ParsedTypeLine) IsBasic() bool {
	return contains(p.Supertypes, "BASIC")
}

// HasSubtype returns true if the type line includes the specified subtype
func (p ParsedTypeLine) HasSubtype(subtype string) bool {
	return contains(p.Subtypes, strings.ToUpper(subtype))
}
