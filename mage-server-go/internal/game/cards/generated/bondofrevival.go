package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Bond Of Revival", NewBondOfRevival)
}

// NewBondOfRevival creates a Bond Of Revival
// {4}{B} - SORCERY
func NewBondOfRevival(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bond Of Revival")
	card.ManaCost = "{4}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect("HasteAbility", effects.DurationUntilYourNextTurn)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
