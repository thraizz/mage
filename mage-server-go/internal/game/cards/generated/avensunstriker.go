package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Aven Sunstriker", NewAvenSunstriker)
}

// NewAvenSunstriker creates a Aven Sunstriker
// {1}{W}{W} - CREATURE
// Flying, DoubleStrike
func NewAvenSunstriker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Aven Sunstriker")
	card.ManaCost = "{1}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BIRD", "WARRIOR"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDoubleStrike)
	card.AddAbility(ability1)
	return card, nil
}
