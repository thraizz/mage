package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("You Look Upon The Tarrasque", NewYouLookUponTheTarrasque)
}

// NewYouLookUponTheTarrasque creates a You Look Upon The Tarrasque
// {4}{G} - INSTANT
func NewYouLookUponTheTarrasque(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "You Look Upon The Tarrasque")
	card.ManaCost = "{4}{G}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(5, 5)).
		AddEffect(abilities.NewGrantAbilityEffect("IndestructibleAbility")).
		AddTarget(abilities.NewCreatureTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}