package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Earthbender Ascension", NewEarthbenderAscension)
}

// NewEarthbenderAscension creates a Earthbender Ascension
// {2}{G} - ENCHANTMENT
func NewEarthbenderAscension(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Earthbender Ascension")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EntersBattlefieldTriggeredAbility
	//   - Effect: EarthbendTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
