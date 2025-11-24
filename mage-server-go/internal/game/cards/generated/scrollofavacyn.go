package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Scroll Of Avacyn", NewScrollOfAvacyn)
}

// NewScrollOfAvacyn creates a Scroll Of Avacyn
// {1} - ARTIFACT
func NewScrollOfAvacyn(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Scroll Of Avacyn")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{1}").
		AddSacrificeSourceCost().
		AddEffect(abilities.NewDrawCardsEffect(1)).
		AddEffect(abilities.NewGainLifeEffect(5)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}