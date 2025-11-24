package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Undergrowth Recon", NewUndergrowthRecon)
}

// NewUndergrowthRecon creates a Undergrowth Recon
// {1}{G}{G} - ENCHANTMENT
func NewUndergrowthRecon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Undergrowth Recon")
	card.ManaCost = "{1}{G}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
