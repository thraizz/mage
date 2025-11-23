package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Feywild Visitor", NewFeywildVisitor)
}

// NewFeywildVisitor creates a Feywild Visitor
// {2}{U} - ENCHANTMENT
func NewFeywildVisitor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Feywild Visitor")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"BACKGROUND"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("FaerieDragonToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token0_0)).
		AddEffect(abilities.NewGrantAbilityEffect(new OneOrMoreCombatDamagePlayerTriggeredAbility(new CreateTokenEffect(token0_0), ) .setTriggerPhrase("Whenever one or more nontoken creatures you control deal combat damage to a player, you "))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}