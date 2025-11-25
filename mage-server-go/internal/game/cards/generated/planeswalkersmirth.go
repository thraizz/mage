package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Planeswalkers Mirth", NewPlaneswalkersMirth)
}

// NewPlaneswalkersMirth creates a Planeswalkers Mirth
// {2}{W} - ENCHANTMENT
func NewPlaneswalkersMirth(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Planeswalkers Mirth")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - PlaneswalkersMirthEffect()
	// card.AddAbility(ability0)
	return card, nil
}
