package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Archmages Charm", NewArchmagesCharm)
}

// NewArchmagesCharm creates a Archmages Charm
// {U}{U}{U} - INSTANT
func NewArchmagesCharm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Archmages Charm")
	card.ManaCost = "{U}{U}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(2)).
		AddEffect(abilities.NewGainControlTargetEffect(abilities.DurationCustom)).
		// TODO: CounterTargetEffect with complex parameters
		AddTarget(abilities.NewSpellTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
