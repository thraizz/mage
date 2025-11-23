package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jade Idol", NewJadeIdol)
}

// NewJadeIdol creates a Jade Idol
// {4} - ARTIFACT
func NewJadeIdol(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jade Idol")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"SPIRIT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
