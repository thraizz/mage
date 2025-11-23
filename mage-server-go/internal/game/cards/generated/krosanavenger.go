package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Krosan Avenger", NewKrosanAvenger)
}

// NewKrosanAvenger creates a Krosan Avenger
// {2}{G} - CREATURE
// Trample
func NewKrosanAvenger(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Krosan Avenger")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "DRUID"}
	card.Power = "3"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - RegenerateSourceEffect()
	// card.AddAbility(ability1)
	return card, nil
}
