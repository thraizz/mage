package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Knotvine Mystic", NewKnotvineMystic)
}

// NewKnotvineMystic creates a Knotvine Mystic
// {R}{G}{W} - CREATURE
func NewKnotvineMystic(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Knotvine Mystic")
	card.ManaCost = "{R}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "DRUID"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
