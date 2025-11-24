package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sanctum Of Shattered Heights", NewSanctumOfShatteredHeights)
}

// NewSanctumOfShatteredHeights creates a Sanctum Of Shattered Heights
// {2}{R} - ENCHANTMENT
func NewSanctumOfShatteredHeights(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sanctum Of Shattered Heights")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SHRINE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewDamageEffect(xValue)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
