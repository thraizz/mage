package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Apex Devastator", NewApexDevastator)
}

// NewApexDevastator creates a Apex Devastator
// {8}{G}{G} - CREATURE
func NewApexDevastator(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Apex Devastator")
	card.ManaCost = "{8}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"CHIMERA", "HYDRA"}
	card.Power = "10"
	card.Toughness = "10"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
