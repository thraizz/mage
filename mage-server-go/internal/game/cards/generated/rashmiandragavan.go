package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rashmi And Ragavan", NewRashmiAndRagavan)
}

// NewRashmiAndRagavan creates a Rashmi And Ragavan
// {1}{G}{U}{R} - CREATURE
func NewRashmiAndRagavan(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rashmi And Ragavan")
	card.ManaCost = "{1}{G}{U}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "MONKEY"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
