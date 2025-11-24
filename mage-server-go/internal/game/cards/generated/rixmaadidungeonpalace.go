package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rix Maadi Dungeon Palace", NewRixMaadiDungeonPalace)
}

// NewRixMaadiDungeonPalace creates a Rix Maadi Dungeon Palace
//   - LAND
func NewRixMaadiDungeonPalace(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rix Maadi Dungeon Palace")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DiscardEachPlayerEffect()
	// card.AddAbility(ability1)
	return card, nil
}
