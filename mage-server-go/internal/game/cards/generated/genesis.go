package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Genesis", NewGenesis)
}

// NewGenesis creates a Genesis
// {4}{G} - CREATURE
func NewGenesis(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Genesis")
	card.ManaCost = "{4}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"INCARNATION"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: BeginningOfUpkeepTriggeredAbility
	//   - Effect: DoIfCostPaid(new ReturnFromGraveyardToHandTargetEffect(), new M...)
	// card.AddAbility(ability0)
	return card, nil
}
