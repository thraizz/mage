package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Suffocating Blast", NewSuffocatingBlast)
}

// NewSuffocatingBlast creates a Suffocating Blast
// {1}{U}{U}{R} - INSTANT
func NewSuffocatingBlast(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Suffocating Blast")
	card.ManaCost = "{1}{U}{U}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(3)).
		// TODO: CounterTargetEffect with complex parameters
		AddEffect(abilities.NewDamageEffect(3)).
		AddTarget(abilities.NewSpellTargetFilter()).
		AddTarget(abilities.NewCreatureTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
