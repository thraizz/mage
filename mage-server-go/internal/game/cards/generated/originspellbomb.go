package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Origin Spellbomb", NewOriginSpellbomb)
}

// NewOriginSpellbomb creates a Origin Spellbomb
// {1} - ARTIFACT
func NewOriginSpellbomb(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Origin Spellbomb")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("MyrToken")
	if err != nil {
		return nil, err
	}
	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{1}").
		AddTapCost().
		AddSacrificeSourceCost().
		AddEffect(abilities.NewCreateTokenEffect(token0_0, 1)).
		Build()
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new DrawCardSourceControllerEffect(1), new ManaCos...)
	// card.AddAbility(ability1)
	return card, nil
}
