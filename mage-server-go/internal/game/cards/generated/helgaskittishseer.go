package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Helga Skittish Seer", NewHelgaSkittishSeer)
}

// NewHelgaSkittishSeer creates a Helga Skittish Seer
// {G}{W}{U} - CREATURE
func NewHelgaSkittishSeer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Helga Skittish Seer")
	card.ManaCost = "{G}{W}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"FROG", "DRUID"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
