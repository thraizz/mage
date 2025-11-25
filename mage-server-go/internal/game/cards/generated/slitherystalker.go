package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Slithery Stalker", NewSlitheryStalker)
}

// NewSlitheryStalker creates a Slithery Stalker
// {1}{B}{B} - CREATURE
func NewSlitheryStalker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Slithery Stalker")
	card.ManaCost = "{1}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"NIGHTMARE", "HORROR"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EntersBattlefieldTriggeredAbility
	//   - Effect: ExileTargetForSourceEffect()
	// card.AddAbility(ability0)
	// TODO: Implement triggered ability: LeavesBattlefieldTriggeredAbility
	//   - Effect: ReturnFromExileForSourceEffect()
	// card.AddAbility(ability1)
	return card, nil
}
