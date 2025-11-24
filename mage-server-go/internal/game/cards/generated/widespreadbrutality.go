package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Widespread Brutality", NewWidespreadBrutality)
}

// NewWidespreadBrutality creates a Widespread Brutality
// {1}{B}{R}{R} - SORCERY
func NewWidespreadBrutality(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Widespread Brutality")
	card.ManaCost = "{1}{B}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}