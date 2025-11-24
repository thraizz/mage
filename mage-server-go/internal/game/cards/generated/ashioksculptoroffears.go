package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ashiok Sculptor Of Fears", NewAshiokSculptorOfFears)
}

// NewAshiokSculptorOfFears creates a Ashiok Sculptor Of Fears
// {4}{U}{B} - PLANESWALKER
func NewAshiokSculptorOfFears(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ashiok Sculptor Of Fears")
	card.ManaCost = "{4}{U}{B}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"ASHIOK"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - ReturnFromGraveyardToBattlefieldTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}