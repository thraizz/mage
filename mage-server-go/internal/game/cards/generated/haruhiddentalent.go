package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Haru Hidden Talent", NewHaruHiddenTalent)
}

// NewHaruHiddenTalent creates a Haru Hidden Talent
// {1}{G} - CREATURE
func NewHaruHiddenTalent(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Haru Hidden Talent")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "PEASANT", "ALLY"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EntersBattlefieldAllTriggeredAbility
	//   - Effect: EarthbendTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
