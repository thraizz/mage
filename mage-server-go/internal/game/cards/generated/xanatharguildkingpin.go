package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Xanathar Guild Kingpin", NewXanatharGuildKingpin)
}

// NewXanatharGuildKingpin creates a Xanathar Guild Kingpin
// {4}{U}{B} - CREATURE
func NewXanatharGuildKingpin(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Xanathar Guild Kingpin")
	card.ManaCost = "{4}{U}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BEHOLDER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: BeginningOfUpkeepTriggeredAbility
	//   - Effect: XanatharGuildKingpinRuleModifyingEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewOpponentTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
