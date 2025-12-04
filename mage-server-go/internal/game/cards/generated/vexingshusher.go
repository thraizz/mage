package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vexing Shusher", NewVexingShusher)
}

// NewVexingShusher creates a Vexing Shusher
// {R/G}{R/G} - CREATURE
func NewVexingShusher(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vexing Shusher")
	card.ManaCost = "{R/G}{R/G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GOBLIN", "SHAMAN"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - VexingShusherCantCounterTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
