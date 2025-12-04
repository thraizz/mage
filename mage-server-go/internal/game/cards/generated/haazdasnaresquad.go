package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Haazda Snare Squad", NewHaazdaSnareSquad)
}

// NewHaazdaSnareSquad creates a Haazda Snare Squad
// {2}{W} - CREATURE
func NewHaazdaSnareSquad(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Haazda Snare Squad")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "SOLDIER"}
	card.Power = "1"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: AttacksTriggeredAbility
	//   - Effect: DoIfCostPaid(new TapTargetEffect(), new ManaCostsImpl<>("{W}"))
	// card.AddAbility(ability0)
	return card, nil
}
