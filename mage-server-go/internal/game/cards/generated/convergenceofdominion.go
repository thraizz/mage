package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Convergence Of Dominion", NewConvergenceOfDominion)
}

// NewConvergenceOfDominion creates a Convergence Of Dominion
// {3} - ARTIFACT
func NewConvergenceOfDominion(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Convergence Of Dominion")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{3}").
		AddTapCost().
		AddEffect(abilities.NewMillCardsControllerEffect(1)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}