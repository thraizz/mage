package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sands Of Delirium", NewSandsOfDelirium)
}

// NewSandsOfDelirium creates a Sands Of Delirium
// {3} - ARTIFACT
func NewSandsOfDelirium(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sands Of Delirium")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewMillCardsTargetEffect(1)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
