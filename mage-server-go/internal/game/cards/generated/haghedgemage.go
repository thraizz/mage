package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hag Hedge Mage", NewHagHedgeMage)
}

// NewHagHedgeMage creates a Hag Hedge Mage
// {2}{B/G} - CREATURE
func NewHagHedgeMage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hag Hedge Mage")
	card.ManaCost = "{2}{B/G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HAG", "SHAMAN"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EntersBattlefieldTriggeredAbility
	//   - Effect: DiscardTargetEffect(1)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPlayerTargetFilter())
	// card.AddAbility(ability0)
	// TODO: Implement triggered ability: EntersBattlefieldTriggeredAbility
	//   - Effect: PutOnLibraryTargetEffect()
	// card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(1)
	// card.AddAbility(ability2)
	return card, nil
}
