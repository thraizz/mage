package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Beguiler Of Wills", NewBeguilerOfWills)
}

// NewBeguilerOfWills creates a Beguiler Of Wills
// {3}{U}{U} - CREATURE
func NewBeguilerOfWills(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Beguiler Of Wills")
	card.ManaCost = "{3}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewGainControlTargetEffect(abilities.DurationCustom)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}