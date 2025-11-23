package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bloodtallow Candle", NewBloodtallowCandle)
}

// NewBloodtallowCandle creates a Bloodtallow Candle
// {1} - ARTIFACT
func NewBloodtallowCandle(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bloodtallow Candle")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{6}").
		AddTapCost().
		AddSacrificeSourceCost().
		AddEffect(abilities.NewBoostEffect(-5, -5)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
