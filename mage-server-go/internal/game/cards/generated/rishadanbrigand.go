package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rishadan Brigand", NewRishadanBrigand)
}

// NewRishadanBrigand creates a Rishadan Brigand
// {4}{U} - CREATURE
// Flying
func NewRishadanBrigand(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rishadan Brigand")
	card.ManaCost = "{4}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "PIRATE"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
