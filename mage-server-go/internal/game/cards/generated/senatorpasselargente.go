package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Senator Passel Argente", NewSenatorPasselArgente)
}

// NewSenatorPasselArgente creates a Senator Passel Argente
// {B}{R} - CREATURE
func NewSenatorPasselArgente(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Senator Passel Argente")
	card.ManaCost = "{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KOORIVAR", "ADVISOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
