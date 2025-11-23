package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Needlebite Trap", NewNeedlebiteTrap)
}

// NewNeedlebiteTrap creates a Needlebite Trap
// {5}{B}{B} - INSTANT
func NewNeedlebiteTrap(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Needlebite Trap")
	card.ManaCost = "{5}{B}{B}"
	card.Types = []string{"INSTANT"}
	card.Subtypes = []string{"TRAP"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewLoseLifeEffect(5)).
		AddEffect(abilities.NewGainLifeEffect(5)).
		AddTarget(abilities.NewPlayerTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
