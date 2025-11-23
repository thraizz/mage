package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Avenging Angel", NewAvengingAngel)
}

// NewAvengingAngel creates a Avenging Angel
// {3}{W}{W} - CREATURE
// Flying
func NewAvengingAngel(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Avenging Angel")
	card.ManaCost = "{3}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ANGEL"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
