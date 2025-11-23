package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Phyrexian Juggernaut", NewPhyrexianJuggernaut)
}

// NewPhyrexianJuggernaut creates a Phyrexian Juggernaut
// {6} - ARTIFACT CREATURE
func NewPhyrexianJuggernaut(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Phyrexian Juggernaut")
	card.ManaCost = "{6}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "JUGGERNAUT"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
