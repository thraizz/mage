package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Goreclaw Terror Of Qal Sisma", NewGoreclawTerrorOfQalSisma)
}

// NewGoreclawTerrorOfQalSisma creates a Goreclaw Terror Of Qal Sisma
// {3}{G} - CREATURE
func NewGoreclawTerrorOfQalSisma(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Goreclaw Terror Of Qal Sisma")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BEAR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
