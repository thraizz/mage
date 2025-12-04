package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hisoka Minamo Sensei", NewHisokaMinamoSensei)
}

// NewHisokaMinamoSensei creates a Hisoka Minamo Sensei
// {2}{U}{U} - CREATURE
func NewHisokaMinamoSensei(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hisoka Minamo Sensei")
	card.ManaCost = "{2}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - HisokaMinamoSenseiCounterEffect()
	// card.AddAbility(ability0)
	return card, nil
}
