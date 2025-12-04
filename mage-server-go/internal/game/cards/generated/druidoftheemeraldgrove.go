package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Druid Of The Emerald Grove", NewDruidOfTheEmeraldGrove)
}

// NewDruidOfTheEmeraldGrove creates a Druid Of The Emerald Grove
// {3}{G} - CREATURE
func NewDruidOfTheEmeraldGrove(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Druid Of The Emerald Grove")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DWARF", "DRUID"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
