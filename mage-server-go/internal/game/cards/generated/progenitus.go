package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Progenitus", NewProgenitus)
}

// NewProgenitus creates a Progenitus
// {W}{W}{U}{U}{B}{B}{R}{R}{G}{G} - CREATURE
func NewProgenitus(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Progenitus")
	card.ManaCost = "{W}{W}{U}{U}{B}{B}{R}{R}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HYDRA", "AVATAR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "10"
	card.Toughness = "10"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
