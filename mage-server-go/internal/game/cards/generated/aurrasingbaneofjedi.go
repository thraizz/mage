package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Aurra Sing Bane Of Jedi", NewAurraSingBaneOfJedi)
}

// NewAurraSingBaneOfJedi creates a Aurra Sing Bane Of Jedi
// {2}{B}{R} - PLANESWALKER
func NewAurraSingBaneOfJedi(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Aurra Sing Bane Of Jedi")
	card.ManaCost = "{2}{B}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"AURRA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
