package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Skirk Drill Sergeant", NewSkirkDrillSergeant)
}

// NewSkirkDrillSergeant creates a Skirk Drill Sergeant
// {1}{R} - CREATURE
func NewSkirkDrillSergeant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Skirk Drill Sergeant")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GOBLIN"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new SkirkDrillSergeantEffect(), new ManaCostsImpl<...)
	// card.AddAbility(ability0)
	return card, nil
}