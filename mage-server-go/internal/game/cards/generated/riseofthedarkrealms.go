package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rise Of The Dark Realms", NewRiseOfTheDarkRealms)
}

// NewRiseOfTheDarkRealms creates a Rise Of The Dark Realms
// {7}{B}{B} - SORCERY
func NewRiseOfTheDarkRealms(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rise Of The Dark Realms")
	card.ManaCost = "{7}{B}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}