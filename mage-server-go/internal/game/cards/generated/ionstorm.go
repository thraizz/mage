package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ion Storm", NewIonStorm)
}

// NewIonStorm creates a Ion Storm
// {2}{R} - ENCHANTMENT
func NewIonStorm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ion Storm")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewDamageEffect(2)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
