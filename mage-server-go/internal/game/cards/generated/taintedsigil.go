package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tainted Sigil", NewTaintedSigil)
}

// NewTaintedSigil creates a Tainted Sigil
// {1}{W}{B} - ARTIFACT
func NewTaintedSigil(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tainted Sigil")
	card.ManaCost = "{1}{W}{B}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddSacrificeSourceCost().
		AddEffect(abilities.NewGainLifeEffect(AllPlayersLostLifeCount.instance, rule)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}