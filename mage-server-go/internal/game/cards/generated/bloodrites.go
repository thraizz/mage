package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Blood Rites", NewBloodRites)
}

// NewBloodRites creates a Blood Rites
// {3}{R}{R} - ENCHANTMENT
func NewBloodRites(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Blood Rites")
	card.ManaCost = "{3}{R}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewDamageEffect(2)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
