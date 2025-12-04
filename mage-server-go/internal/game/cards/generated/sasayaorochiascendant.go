package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sasaya Orochi Ascendant", NewSasayaOrochiAscendant)
}

// NewSasayaOrochiAscendant creates a Sasaya Orochi Ascendant
// {1}{G}{G} - CREATURE
func NewSasayaOrochiAscendant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sasaya Orochi Ascendant")
	card.ManaCost = "{1}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SNAKE", "MONK"}
	card.Supertypes = []string{"LEGENDARY", "LEGENDARY"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
