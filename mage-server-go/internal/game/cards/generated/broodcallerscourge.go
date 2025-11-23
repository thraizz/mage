package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Broodcaller Scourge", NewBroodcallerScourge)
}

// NewBroodcallerScourge creates a Broodcaller Scourge
// {5}{G}{G} - CREATURE
// Flying
func NewBroodcallerScourge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Broodcaller Scourge")
	card.ManaCost = "{5}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DRAGON"}
	card.Power = "5"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
