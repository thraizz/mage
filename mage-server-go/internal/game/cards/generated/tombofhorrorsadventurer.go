package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tomb Of Horrors Adventurer", NewTombOfHorrorsAdventurer)
}

// NewTombOfHorrorsAdventurer creates a Tomb Of Horrors Adventurer
// {5}{U} - CREATURE
func NewTombOfHorrorsAdventurer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tomb Of Horrors Adventurer")
	card.ManaCost = "{5}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "MONK"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
