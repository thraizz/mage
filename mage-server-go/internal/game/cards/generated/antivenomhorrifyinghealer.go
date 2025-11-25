package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Anti Venom Horrifying Healer", NewAntiVenomHorrifyingHealer)
}

// NewAntiVenomHorrifyingHealer creates a Anti Venom Horrifying Healer
// {W}{W}{W}{W}{W} - CREATURE
func NewAntiVenomHorrifyingHealer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Anti Venom Horrifying Healer")
	card.ManaCost = "{W}{W}{W}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SYMBIOTE", "HERO"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EntersBattlefieldTriggeredAbility
	//   - Effect: ReturnFromGraveyardToBattlefieldTargetEffect()
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - ReturnFromGraveyardToBattlefieldTargetEffect()
	// card.AddAbility(ability1)
	return card, nil
}
