package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ulamogs Despoiler", NewUlamogsDespoiler)
}

// NewUlamogsDespoiler creates a Ulamogs Despoiler
// {6} - CREATURE
func NewUlamogsDespoiler(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ulamogs Despoiler")
	card.ManaCost = "{6}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDRAZI", "PROCESSOR"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}