package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Spellbinding Soprano", NewSpellbindingSoprano)
}

// NewSpellbindingSoprano creates a Spellbinding Soprano
// {1}{R} - CREATURE
func NewSpellbindingSoprano(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Spellbinding Soprano")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "BARD"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
