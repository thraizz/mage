package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Melded Moxite", NewMeldedMoxite)
}

// NewMeldedMoxite creates a Melded Moxite
// {1}{R} - ARTIFACT
func NewMeldedMoxite(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Melded Moxite")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("RobotToken")
	if err != nil {
		return nil, err
	}
	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{3}").
		AddSacrificeSourceCost().
		AddEffect(abilities.NewCreateTokenEffectTapped(token0_0, 1, true)).
		Build()
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new DrawCardSourceControllerEffect(2), new Discard...)
	// card.AddAbility(ability1)
	return card, nil
}
