package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rise Of The Witch King", NewRiseOfTheWitchKing)
}

// NewRiseOfTheWitchKing creates a Rise Of The Witch King
// {2}{B}{G} - SORCERY
func NewRiseOfTheWitchKing(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rise Of The Witch King")
	card.ManaCost = "{2}{B}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
