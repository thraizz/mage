package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kathril Aspect Warper", NewKathrilAspectWarper)
}

// NewKathrilAspectWarper creates a Kathril Aspect Warper
// {2}{W}{B}{G} - CREATURE
func NewKathrilAspectWarper(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kathril Aspect Warper")
	card.ManaCost = "{2}{W}{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"NIGHTMARE", "INSECT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}