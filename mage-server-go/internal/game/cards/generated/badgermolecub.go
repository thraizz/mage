package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Badgermole Cub", NewBadgermoleCub)
}

// NewBadgermoleCub creates a Badgermole Cub
// {1}{G} - CREATURE
func NewBadgermoleCub(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Badgermole Cub")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BADGER", "MOLE"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EntersBattlefieldTriggeredAbility
	//   - Effect: EarthbendTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
