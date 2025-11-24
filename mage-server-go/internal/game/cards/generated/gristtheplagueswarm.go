package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Grist The Plague Swarm", NewGristThePlagueSwarm)
}

// NewGristThePlagueSwarm creates a Grist The Plague Swarm
//   - PLANESWALKER
func NewGristThePlagueSwarm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Grist The Plague Swarm")
	card.ManaCost = ""
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"GRIST"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("IzoniInsectToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: DestroyTargetEffect with complex parameters
		AddEffect(abilities.NewCreateTokenEffect(token0_0)).
		AddEffect(abilities.NewAddCountersTargetEffect(counters.NewCounter("deathtouch", 1))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
