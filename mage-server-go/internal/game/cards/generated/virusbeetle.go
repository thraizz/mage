package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Virus Beetle", NewVirusBeetle)
}

// NewVirusBeetle creates a Virus Beetle
// {1}{B} - ARTIFACT CREATURE
func NewVirusBeetle(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Virus Beetle")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"INSECT"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardEachPlayerEffect(TargetController.OPPONENT)
	// card.AddAbility(ability0)
	return card, nil
}
