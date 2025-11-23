package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Chalice Of The Void", NewChaliceOfTheVoid)
}

// NewChaliceOfTheVoid creates a Chalice Of The Void
// {X}{X} - ARTIFACT
func NewChaliceOfTheVoid(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Chalice Of The Void")
	card.ManaCost = "{X}{X}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
