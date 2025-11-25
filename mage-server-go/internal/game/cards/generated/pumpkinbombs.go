package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Pumpkin Bombs", NewPumpkinBombs)
}

// NewPumpkinBombs creates a Pumpkin Bombs
// {1}{R} - ARTIFACT
func NewPumpkinBombs(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Pumpkin Bombs")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewDrawCardsEffect(1)).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.NewCounter("fuse", 1))).
		AddEffect(abilities.NewDamageEffect(xValue)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
