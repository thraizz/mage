package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Pursuit Of Knowledge", NewPursuitOfKnowledge)
}

// NewPursuitOfKnowledge creates a Pursuit Of Knowledge
// {3}{W} - ENCHANTMENT
func NewPursuitOfKnowledge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Pursuit Of Knowledge")
	card.ManaCost = "{3}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddSacrificeSourceCost().
		AddEffect(abilities.NewDrawCardsEffect(7)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
