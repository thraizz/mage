package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Death Cloud", NewDeathCloud)
}

// NewDeathCloud creates a Death Cloud
// {X}{B}{B}{B} - SORCERY
func NewDeathCloud(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Death Cloud")
	card.ManaCost = "{X}{B}{B}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeAllEffect(xValue, new FilterControlledCreaturePermanent("cre...)
	//   - DiscardEachPlayerEffect(xValue, false)
	// card.AddAbility(ability0)
	return card, nil
}
