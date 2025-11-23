package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Diligent Excavator", NewDiligentExcavator)
}

// NewDiligentExcavator creates a Diligent Excavator
// {1}{U} - CREATURE
func NewDiligentExcavator(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Diligent Excavator")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ARTIFICER"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewMillCardsTargetEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
