package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Ragavan Nimble Pilferer", NewRagavanNimblePilferer)
}

// NewRagavanNimblePilferer creates a Ragavan Nimble Pilferer
// {R} - CREATURE
func NewRagavanNimblePilferer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ragavan Nimble Pilferer")
	card.ManaCost = "{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MONKEY", "PIRATE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("TreasureToken")
	if err != nil {
		return nil, err
	}
	ability0 := abilities.NewTriggeredAbilityBuilder(card.ID).
		// TODO: Set trigger for LeavesBattlefieldAll (when any permanent you control leaves the battlefield)
		// SetTrigger(abilities.NewLeavesBattlefieldAllTrigger(card.ID, abilities.NewControlledPermanentFilter())).
		AddEffect(abilities.NewCreateTokenEffect(token0_0)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
