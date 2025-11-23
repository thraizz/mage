package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Invasion Of Dominaria", NewInvasionOfDominaria)
}

// NewInvasionOfDominaria creates a Invasion Of Dominaria
// {2}{W} - BATTLE
func NewInvasionOfDominaria(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Invasion Of Dominaria")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"BATTLE"}
	card.Subtypes = []string{"SIEGE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainLifeEffect(4)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
