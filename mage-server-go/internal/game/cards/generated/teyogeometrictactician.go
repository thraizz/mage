package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Teyo Geometric Tactician", NewTeyoGeometricTactician)
}

// NewTeyoGeometricTactician creates a Teyo Geometric Tactician
// {2}{W} - PLANESWALKER
func NewTeyoGeometricTactician(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Teyo Geometric Tactician")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"TEYO"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewTriggeredAbilityBuilder(card.ID).
		// TODO: Set trigger for LeavesBattlefieldAll (when any permanent you control leaves the battlefield)
		// SetTrigger(abilities.NewLeavesBattlefieldAllTrigger(card.ID, abilities.NewControlledPermanentFilter())).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		AddTarget(abilities.NewOpponentTargetFilter()).
		Build()
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("WallFlyingToken")
	if err != nil {
		return nil, err
	}
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token1_0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
