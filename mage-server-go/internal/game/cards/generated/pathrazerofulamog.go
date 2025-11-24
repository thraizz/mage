package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Pathrazer Of Ulamog", NewPathrazerOfUlamog)
}

// NewPathrazerOfUlamog creates a Pathrazer Of Ulamog
// {11} - CREATURE
func NewPathrazerOfUlamog(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Pathrazer Of Ulamog")
	card.ManaCost = "{11}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDRAZI"}
	card.Power = "9"
	card.Toughness = "9"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
