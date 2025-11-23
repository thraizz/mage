package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Choice Of Damnations", NewChoiceOfDamnations)
}

// NewChoiceOfDamnations creates a Choice Of Damnations
// {5}{B} - SORCERY
func NewChoiceOfDamnations(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Choice Of Damnations")
	card.ManaCost = "{5}{B}"
	card.Types = []string{"SORCERY"}
	card.Subtypes = []string{"ARCANE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddTarget(abilities.NewOpponentTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
