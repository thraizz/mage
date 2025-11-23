package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Endless Cockroaches", NewEndlessCockroaches)
}

// NewEndlessCockroaches creates a Endless Cockroaches
// {1}{B}{B} - CREATURE
func NewEndlessCockroaches(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Endless Cockroaches")
	card.ManaCost = "{1}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"INSECT"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewReturnToHandSourceEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
