package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Trevas Ruins", NewTrevasRuins)
}

// NewTrevasRuins creates a Trevas Ruins
//   - LAND
func NewTrevasRuins(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Trevas Ruins")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"LAIR"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "G")
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "W")
	card.AddAbility(ability1)
	ability2 := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability2)
	return card, nil
}
