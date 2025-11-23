package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Chandra Torch Of Defiance", NewChandraTorchOfDefiance)
}

// NewChandraTorchOfDefiance creates a Chandra Torch Of Defiance
// {2}{R}{R} - PLANESWALKER
func NewChandraTorchOfDefiance(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Chandra Torch Of Defiance")
	card.ManaCost = "{2}{R}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"CHANDRA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(4)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
