package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kheru Spellsnatcher", NewKheruSpellsnatcher)
}

// NewKheruSpellsnatcher creates a Kheru Spellsnatcher
// {3}{U} - CREATURE
func NewKheruSpellsnatcher(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kheru Spellsnatcher")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SNAKE", "WIZARD"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: TurnedFaceUpSourceTriggeredAbility
	//   - Effect: KheruSpellsnatcherEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewSpellTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
