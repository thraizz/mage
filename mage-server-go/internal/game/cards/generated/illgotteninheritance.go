package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ill Gotten Inheritance", NewIllGottenInheritance)
}

// NewIllGottenInheritance creates a Ill Gotten Inheritance
// {3}{B} - ENCHANTMENT
func NewIllGottenInheritance(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ill Gotten Inheritance")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddSacrificeSourceCost().
		AddEffect(abilities.NewDamageEffect(4)).
		AddEffect(abilities.NewGainLifeEffect(4)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}