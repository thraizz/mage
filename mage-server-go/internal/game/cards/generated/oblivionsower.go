package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Oblivion Sower", NewOblivionSower)
}

// NewOblivionSower creates a Oblivion Sower
// {6} - CREATURE
func NewOblivionSower(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Oblivion Sower")
	card.ManaCost = "{6}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDRAZI"}
	card.Power = "5"
	card.Toughness = "8"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
