package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Invasion Of Ravnica", NewInvasionOfRavnica)
}

// NewInvasionOfRavnica creates a Invasion Of Ravnica
// {5} - BATTLE
func NewInvasionOfRavnica(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Invasion Of Ravnica")
	card.ManaCost = "{5}"
	card.Types = []string{"BATTLE"}
	card.Subtypes = []string{"SIEGE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: ExileTargetEffect with complex parameters
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
