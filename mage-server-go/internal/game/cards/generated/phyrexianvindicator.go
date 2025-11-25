package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Phyrexian Vindicator", NewPhyrexianVindicator)
}

// NewPhyrexianVindicator creates a Phyrexian Vindicator
// {W}{W}{W}{W} - CREATURE
// Flying
func NewPhyrexianVindicator(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Phyrexian Vindicator")
	card.ManaCost = "{W}{W}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
