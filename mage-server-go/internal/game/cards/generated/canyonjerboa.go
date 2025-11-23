package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Canyon Jerboa", NewCanyonJerboa)
}

// NewCanyonJerboa creates a Canyon Jerboa
// {2}{W} - CREATURE
func NewCanyonJerboa(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Canyon Jerboa")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MOUSE"}
	card.Power = "1"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(1, 1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
