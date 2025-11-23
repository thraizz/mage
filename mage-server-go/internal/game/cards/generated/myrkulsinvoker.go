package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Myrkuls Invoker", NewMyrkulsInvoker)
}

// NewMyrkulsInvoker creates a Myrkuls Invoker
// {2}{B} - CREATURE
func NewMyrkulsInvoker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Myrkuls Invoker")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"CAT", "ROGUE"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{8}").
		AddEffect(abilities.NewBoostEffect(2, 0)).
		AddEffect(abilities.NewGrantAbilityEffect(new MenaceAbility(false))).
		Build()
	card.AddAbility(ability0)
	return card, nil
}