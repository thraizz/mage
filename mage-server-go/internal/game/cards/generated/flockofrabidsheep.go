package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Flock Of Rabid Sheep", NewFlockOfRabidSheep)
}

// NewFlockOfRabidSheep creates a Flock Of Rabid Sheep
// {X}{G}{G} - SORCERY
func NewFlockOfRabidSheep(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Flock Of Rabid Sheep")
	card.ManaCost = "{X}{G}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
