package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Keldon Battlewagon", NewKeldonBattlewagon)
}

// NewKeldonBattlewagon creates a Keldon Battlewagon
// {5} - ARTIFACT CREATURE
// Trample
func NewKeldonBattlewagon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Keldon Battlewagon")
	card.ManaCost = "{5}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"JUGGERNAUT"}
	card.Power = "0"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeSourceEffect()
	// card.AddAbility(ability1)
	return card, nil
}
