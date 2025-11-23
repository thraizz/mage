package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ghyrson Starn Kelermorph", NewGhyrsonStarnKelermorph)
}

// NewGhyrsonStarnKelermorph creates a Ghyrson Starn Kelermorph
// {1}{U}{R} - CREATURE
func NewGhyrsonStarnKelermorph(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ghyrson Starn Kelermorph")
	card.ManaCost = "{1}{U}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"TYRANID", "HUMAN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
