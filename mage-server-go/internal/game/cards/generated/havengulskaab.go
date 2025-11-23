package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Havengul Skaab", NewHavengulSkaab)
}

// NewHavengulSkaab creates a Havengul Skaab
// {5}{U} - CREATURE
func NewHavengulSkaab(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Havengul Skaab")
	card.ManaCost = "{5}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ZOMBIE", "HORROR"}
	card.Power = "4"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
