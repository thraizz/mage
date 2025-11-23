package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Numa Joraga Chieftain", NewNumaJoragaChieftain)
}

// NewNumaJoragaChieftain creates a Numa Joraga Chieftain
// {2}{G} - CREATURE
func NewNumaJoragaChieftain(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Numa Joraga Chieftain")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "WARRIOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
