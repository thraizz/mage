package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Melek Izzet Paragon", NewMelekIzzetParagon)
}

// NewMelekIzzetParagon creates a Melek Izzet Paragon
// {4}{U}{R} - CREATURE
func NewMelekIzzetParagon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Melek Izzet Paragon")
	card.ManaCost = "{4}{U}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WEIRD", "WIZARD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
