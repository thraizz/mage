package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Candlekeep Sage", NewCandlekeepSage)
}

// NewCandlekeepSage creates a Candlekeep Sage
// {2}{U} - ENCHANTMENT
func NewCandlekeepSage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Candlekeep Sage")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"BACKGROUND"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		AddEffect(abilities.NewGrantAbilityEffect(new EntersBattlefieldOrLeavesSourceTriggeredAbility( new DrawCardSourceControllerEffect(1), false ).setTriggerPhrase("When creature enters or leaves the battlefield, "))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}