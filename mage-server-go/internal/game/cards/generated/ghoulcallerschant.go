package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ghoulcallers Chant", NewGhoulcallersChant)
}

// NewGhoulcallersChant creates a Ghoulcallers Chant
// {B} - SORCERY
func NewGhoulcallersChant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ghoulcallers Chant")
	card.ManaCost = "{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: ReturnFromGraveyardToHandTargetEffect with complex parameters
		// TODO: ReturnFromGraveyardToHandTargetEffect with complex parameters
		// TODO: ReturnFromGraveyardToHandTargetEffect with complex parameters
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
