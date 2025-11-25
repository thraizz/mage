package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Anger Of The Gods", NewAngerOfTheGods)
}

// NewAngerOfTheGods creates a Anger Of The Gods
// {1}{R}{R} - SORCERY
func NewAngerOfTheGods(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Anger Of The Gods")
	card.ManaCost = "{1}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(3, new FilterCreaturePermanent())
	// card.AddAbility(ability0)
	return card, nil
}
