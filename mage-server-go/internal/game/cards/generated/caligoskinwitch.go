package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Caligo Skin Witch", NewCaligoSkinWitch)
}

// NewCaligoSkinWitch creates a Caligo Skin Witch
// {1}{B} - CREATURE
func NewCaligoSkinWitch(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Caligo Skin Witch")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"CREATURE"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardEachPlayerEffect(                 StaticValue.get(2), false, Target...)
	// card.AddAbility(ability0)
	return card, nil
}