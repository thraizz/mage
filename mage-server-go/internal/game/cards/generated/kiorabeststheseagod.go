package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kiora Bests The Sea God", NewKioraBestsTheSeaGod)
}

// NewKioraBestsTheSeaGod creates a Kiora Bests The Sea God
// {5}{U}{U} - ENCHANTMENT
func NewKioraBestsTheSeaGod(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kiora Bests The Sea God")
	card.ManaCost = "{5}{U}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
