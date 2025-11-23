package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vengeful Regrowth", NewVengefulRegrowth)
}

// NewVengefulRegrowth creates a Vengeful Regrowth
// {4}{G}{G} - SORCERY
func NewVengefulRegrowth(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vengeful Regrowth")
	card.ManaCost = "{4}{G}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
