package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ayesha Tanaka Armorer", NewAyeshaTanakaArmorer)
}

// NewAyeshaTanakaArmorer creates a Ayesha Tanaka Armorer
// {3}{W}{U} - CREATURE
func NewAyeshaTanakaArmorer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ayesha Tanaka Armorer")
	card.ManaCost = "{3}{W}{U}"
	card.Types = []string{"CREATURE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 4, Integer.MAX_VALUE, filter, Put...)
	// card.AddAbility(ability0)
	return card, nil
}
