package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Coati Scavenger", NewCoatiScavenger)
}

// NewCoatiScavenger creates a Coati Scavenger
// {2}{G} - CREATURE
func NewCoatiScavenger(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Coati Scavenger")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"RACCOON"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewReturnFromGraveyardToHandTargetEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
