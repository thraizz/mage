package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Savage Twister", NewSavageTwister)
}

// NewSavageTwister creates a Savage Twister
// {X}{R}{G} - SORCERY
func NewSavageTwister(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Savage Twister")
	card.ManaCost = "{X}{R}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(GetXValue.instance, new FilterCreaturePermanent())
	// card.AddAbility(ability0)
	return card, nil
}
