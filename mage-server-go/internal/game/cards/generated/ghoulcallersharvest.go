package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ghoulcallers Harvest", NewGhoulcallersHarvest)
}

// NewGhoulcallersHarvest creates a Ghoulcallers Harvest
// {B}{G} - SORCERY
func NewGhoulcallersHarvest(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ghoulcallers Harvest")
	card.ManaCost = "{B}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
