package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Korlash Heir To Blackblade", NewKorlashHeirToBlackblade)
}

// NewKorlashHeirToBlackblade creates a Korlash Heir To Blackblade
// {2}{B}{B} - CREATURE
func NewKorlashHeirToBlackblade(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Korlash Heir To Blackblade")
	card.ManaCost = "{2}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ZOMBIE", "WARRIOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - RegenerateSourceEffect()
	// card.AddAbility(ability0)
	return card, nil
}
