package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Canopy Vista", NewCanopyVista)
}

// NewCanopyVista creates a Canopy Vista
//   - LAND
func NewCanopyVista(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Canopy Vista")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"FOREST", "PLAINS"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "G")
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "W")
	card.AddAbility(ability1)
	return card, nil
}
