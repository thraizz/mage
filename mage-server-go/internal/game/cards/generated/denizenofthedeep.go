package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Denizen Of The Deep", NewDenizenOfTheDeep)
}

// NewDenizenOfTheDeep creates a Denizen Of The Deep
// {6}{U}{U} - CREATURE
func NewDenizenOfTheDeep(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Denizen Of The Deep")
	card.ManaCost = "{6}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SERPENT"}
	card.Power = "11"
	card.Toughness = "11"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
