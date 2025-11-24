package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Wrath Of The Skies", NewWrathOfTheSkies)
}

// NewWrathOfTheSkies creates a Wrath Of The Skies
// {X}{W}{W} - SORCERY
func NewWrathOfTheSkies(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Wrath Of The Skies")
	card.ManaCost = "{X}{W}{W}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}