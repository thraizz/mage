package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Springjaw Trap", NewSpringjawTrap)
}

// NewSpringjawTrap creates a Springjaw Trap
// {1} - ARTIFACT
// Flash
func NewSpringjawTrap(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Springjaw Trap")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlash)
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{4}").
		AddTapCost().
		AddSacrificeSourceCost().
		AddEffect(abilities.NewDamageEffect(3)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
