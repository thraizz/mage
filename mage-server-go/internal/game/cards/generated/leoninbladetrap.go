package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Leonin Bladetrap", NewLeoninBladetrap)
}

// NewLeoninBladetrap creates a Leonin Bladetrap
// {3} - ARTIFACT
// Flash
func NewLeoninBladetrap(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Leonin Bladetrap")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlash)
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{2}").
		AddSacrificeSourceCost().
		// TODO: DamageAllEffect with complex parameters
		Build()
	card.AddAbility(ability1)
	return card, nil
}
