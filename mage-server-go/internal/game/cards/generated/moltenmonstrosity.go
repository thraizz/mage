package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Molten Monstrosity", NewMoltenMonstrosity)
}

// NewMoltenMonstrosity creates a Molten Monstrosity
// {7}{R} - CREATURE
// Trample
func NewMoltenMonstrosity(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Molten Monstrosity")
	card.ManaCost = "{7}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HELLION"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}