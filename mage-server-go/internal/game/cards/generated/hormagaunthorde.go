package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hormagaunt Horde", NewHormagauntHorde)
}

// NewHormagauntHorde creates a Hormagaunt Horde
// {X}{G} - CREATURE
func NewHormagauntHorde(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hormagaunt Horde")
	card.ManaCost = "{X}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"TYRANID"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                         new ReturnSourceFromGrave...)
	// card.AddAbility(ability0)
	return card, nil
}
