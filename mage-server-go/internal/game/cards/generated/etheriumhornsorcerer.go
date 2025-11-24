package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Etherium Horn Sorcerer", NewEtheriumHornSorcerer)
}

// NewEtheriumHornSorcerer creates a Etherium Horn Sorcerer
// {4}{U}{R} - ARTIFACT CREATURE
func NewEtheriumHornSorcerer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Etherium Horn Sorcerer")
	card.ManaCost = "{4}{U}{R}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"MINOTAUR", "WIZARD"}
	card.Power = "3"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
