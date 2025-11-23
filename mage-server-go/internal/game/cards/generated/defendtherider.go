package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Defend The Rider", NewDefendTheRider)
}

// NewDefendTheRider creates a Defend The Rider
// {G} - INSTANT
func NewDefendTheRider(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Defend The Rider")
	card.ManaCost = "{G}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("PilotSaddleCrewToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect("HexproofAbility")).
		AddEffect(abilities.NewGrantAbilityEffect("IndestructibleAbility")).
		AddEffect(abilities.NewCreateTokenEffect(token0_0)).
		AddTarget(abilities.NewPermanentTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
