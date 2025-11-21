package manual

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

// Register vanilla creatures on package import
func init() {
	cards.Register("Grizzly Bears", NewGrizzlyBears)
	cards.Register("Savannah Lions", NewSavannahLions)
	cards.Register("Hill Giant", NewHillGiant)
}

// NewGrizzlyBears creates a Grizzly Bears
// {1}{G} - Creature — Bear - 2/2
func NewGrizzlyBears(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Grizzly Bears")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BEAR"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M10"
	card.Rarity = "common"
	card.RulesText = ""

	return card, nil
}

// NewSavannahLions creates a Savannah Lions
// {W} - Creature — Cat - 2/1
func NewSavannahLions(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Savannah Lions")
	card.ManaCost = "{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"CAT"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M10"
	card.Rarity = "rare"
	card.RulesText = ""

	return card, nil
}

// NewHillGiant creates a Hill Giant
// {3}{R} - Creature — Giant - 3/3
func NewHillGiant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hill Giant")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GIANT"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M10"
	card.Rarity = "common"
	card.RulesText = ""

	return card, nil
}
