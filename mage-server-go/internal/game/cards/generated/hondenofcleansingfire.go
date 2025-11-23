package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Honden Of Cleansing Fire", NewHondenOfCleansingFire)
}

// NewHondenOfCleansingFire creates a Honden Of Cleansing Fire
// {3}{W} - ENCHANTMENT
func NewHondenOfCleansingFire(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Honden Of Cleansing Fire")
	card.ManaCost = "{3}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SHRINE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainLifeEffect(new MultipliedValue(xValue, 2))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}