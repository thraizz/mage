package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Orchard Warden", NewOrchardWarden)
}

// NewOrchardWarden creates a Orchard Warden
// {4}{G}{G} - CREATURE
func NewOrchardWarden(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Orchard Warden")
	card.ManaCost = "{4}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"TREEFOLK", "SHAMAN"}
	card.Power = "4"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
