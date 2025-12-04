package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Spawnsire Of Ulamog", NewSpawnsireOfUlamog)
}

// NewSpawnsireOfUlamog creates a Spawnsire Of Ulamog
// {10} - CREATURE
func NewSpawnsireOfUlamog(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Spawnsire Of Ulamog")
	card.ManaCost = "{10}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDRAZI"}
	card.Power = "7"
	card.Toughness = "11"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
