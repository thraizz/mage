package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Invasion Of Tolvada", NewInvasionOfTolvada)
}

// NewInvasionOfTolvada creates a Invasion Of Tolvada
// {3}{W}{B} - BATTLE
func NewInvasionOfTolvada(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Invasion Of Tolvada")
	card.ManaCost = "{3}{W}{B}"
	card.Types = []string{"BATTLE"}
	card.Subtypes = []string{"SIEGE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - ReturnFromGraveyardToBattlefieldTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
