package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Devastating Dreams", NewDevastatingDreams)
}

// NewDevastatingDreams creates a Devastating Dreams
// {R}{R} - SORCERY
func NewDevastatingDreams(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Devastating Dreams")
	card.ManaCost = "{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(GetXValue.instance, new FilterCreaturePermanent())
	//   - SacrificeAllEffect(GetXValue.instance, new FilterControlledLandPerman...)
	// card.AddAbility(ability0)
	return card, nil
}
