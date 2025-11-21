package manual

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

// Register all basic lands on package import
func init() {
	cards.Register("Plains", NewPlains)
	cards.Register("Island", NewIsland)
	cards.Register("Swamp", NewSwamp)
	cards.Register("Mountain", NewMountain)
	cards.Register("Forest", NewForest)
}

// NewPlains creates a Plains basic land
// {T}: Add {W}
func NewPlains(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Plains")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"PLAINS"}
	card.Supertypes = []string{"BASIC"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// Add mana ability: {T}: Add {W}
	ability := abilities.BuildSimpleManaAbility(card.ID, "W")
	card.AddAbility(ability)

	return card, nil
}

// NewIsland creates an Island basic land
// {T}: Add {U}
func NewIsland(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Island")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"ISLAND"}
	card.Supertypes = []string{"BASIC"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// Add mana ability: {T}: Add {U}
	ability := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability)

	return card, nil
}

// NewSwamp creates a Swamp basic land
// {T}: Add {B}
func NewSwamp(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Swamp")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"SWAMP"}
	card.Supertypes = []string{"BASIC"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// Add mana ability: {T}: Add {B}
	ability := abilities.BuildSimpleManaAbility(card.ID, "B")
	card.AddAbility(ability)

	return card, nil
}

// NewMountain creates a Mountain basic land
// {T}: Add {R}
func NewMountain(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mountain")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"MOUNTAIN"}
	card.Supertypes = []string{"BASIC"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// Add mana ability: {T}: Add {R}
	ability := abilities.BuildSimpleManaAbility(card.ID, "R")
	card.AddAbility(ability)

	return card, nil
}

// NewForest creates a Forest basic land
// {T}: Add {G}
func NewForest(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Forest")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"FOREST"}
	card.Supertypes = []string{"BASIC"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// Add mana ability: {T}: Add {G}
	ability := abilities.BuildSimpleManaAbility(card.ID, "G")
	card.AddAbility(ability)

	return card, nil
}
