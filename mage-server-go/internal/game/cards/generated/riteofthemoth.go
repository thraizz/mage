package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rite Of The Moth", NewRiteOfTheMoth)
}

// NewRiteOfTheMoth creates a Rite Of The Moth
// {1}{W}{B}{B} - SORCERY
func NewRiteOfTheMoth(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rite Of The Moth")
	card.ManaCost = "{1}{W}{B}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}