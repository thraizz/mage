package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Sibsig Ceremony", NewTheSibsigCeremony)
}

// NewTheSibsigCeremony creates a The Sibsig Ceremony
// {B}{B}{B} - ENCHANTMENT
func NewTheSibsigCeremony(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Sibsig Ceremony")
	card.ManaCost = "{B}{B}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
