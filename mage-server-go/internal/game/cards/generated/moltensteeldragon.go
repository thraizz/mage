package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Moltensteel Dragon", NewMoltensteelDragon)
}

// NewMoltensteelDragon creates a Moltensteel Dragon
// {4}{R/P}{R/P} - ARTIFACT CREATURE
// Flying
func NewMoltensteelDragon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Moltensteel Dragon")
	card.ManaCost = "{4}{R/P}{R/P}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "DRAGON"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
