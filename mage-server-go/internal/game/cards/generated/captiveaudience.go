package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Captive Audience", NewCaptiveAudience)
}

// NewCaptiveAudience creates a Captive Audience
// {5}{B}{R} - ENCHANTMENT
func NewCaptiveAudience(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Captive Audience")
	card.ManaCost = "{5}{B}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: BeginningOfUpkeepTriggeredAbility
	//   - Effect: SetPlayerLifeSourceEffect()
	// card.AddAbility(ability0)
	return card, nil
}
