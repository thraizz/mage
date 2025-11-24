package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Reconnaissance", NewReconnaissance)
}

// NewReconnaissance creates a Reconnaissance
// {W} - ENCHANTMENT
func NewReconnaissance(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Reconnaissance")
	card.ManaCost = "{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{0}").
		AddEffect(abilities.NewUntapEffect("and untap it")).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
