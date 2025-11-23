package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Zhalfirin Knight", NewZhalfirinKnight)
}

// NewZhalfirinKnight creates a Zhalfirin Knight
// {2}{W} - CREATURE
func NewZhalfirinKnight(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Zhalfirin Knight")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"CREATURE"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
