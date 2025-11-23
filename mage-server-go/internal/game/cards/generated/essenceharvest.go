package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Essence Harvest", NewEssenceHarvest)
}

// NewEssenceHarvest creates a Essence Harvest
// {2}{B} - SORCERY
func NewEssenceHarvest(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Essence Harvest")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewLoseLifeEffect(GreatestAmongPermanentsValue.POWER_CONTROLLED_CREATURES)).
		AddEffect(abilities.NewGainLifeEffect(GreatestAmongPermanentsValue.POWER_CONTROLLED_CREATURES)).
		AddTarget(abilities.NewPlayerTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
