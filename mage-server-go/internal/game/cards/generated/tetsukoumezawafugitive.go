package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tetsuko Umezawa Fugitive", NewTetsukoUmezawaFugitive)
}

// NewTetsukoUmezawaFugitive creates a Tetsuko Umezawa Fugitive
// {1}{U} - CREATURE
func NewTetsukoUmezawaFugitive(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tetsuko Umezawa Fugitive")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"CREATURE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
