package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Herd Heirloom", NewHerdHeirloom)
}

// NewHerdHeirloom creates a Herd Heirloom
// {1}{G} - ARTIFACT
func NewHerdHeirloom(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Herd Heirloom")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewGrantAbilityEffect("TrampleAbility")).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		AddEffect(abilities.NewGrantAbilityEffect(new DealsCombatDamageToAPlayerTriggeredAbility(new DrawCardSourceControllerEffect(1)))).
		Build()
	card.AddAbility(ability0)
	return card, nil
}