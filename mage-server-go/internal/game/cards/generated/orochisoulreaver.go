package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Orochi Soul Reaver", NewOrochiSoulReaver)
}

// NewOrochiSoulReaver creates a Orochi Soul Reaver
// {5}{B} - CREATURE
func NewOrochiSoulReaver(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Orochi Soul Reaver")
	card.ManaCost = "{5}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SNAKE", "NINJA", "ROGUE"}
	card.Power = "5"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
