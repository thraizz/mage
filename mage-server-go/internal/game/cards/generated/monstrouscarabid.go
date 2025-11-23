package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Monstrous Carabid", NewMonstrousCarabid)
}

// NewMonstrousCarabid creates a Monstrous Carabid
// {3}{B}{R} - CREATURE
func NewMonstrousCarabid(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Monstrous Carabid")
	card.ManaCost = "{3}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"INSECT"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
