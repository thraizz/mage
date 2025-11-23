package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tectonic Break", NewTectonicBreak)
}

// NewTectonicBreak creates a Tectonic Break
// {X}{R}{R} - SORCERY
func NewTectonicBreak(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tectonic Break")
	card.ManaCost = "{X}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeAllEffect(GetXValue.instance, new FilterControlledLandPerman...)
	// card.AddAbility(ability0)
	return card, nil
}
