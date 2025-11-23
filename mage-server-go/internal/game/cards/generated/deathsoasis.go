package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Deaths Oasis", NewDeathsOasis)
}

// NewDeathsOasis creates a Deaths Oasis
// {W}{B}{G} - ENCHANTMENT
func NewDeathsOasis(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Deaths Oasis")
	card.ManaCost = "{W}{B}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{1}").
		AddSacrificeSourceCost().
		AddEffect(abilities.NewGainLifeEffect(GreatestAmongPermanentsValue.MANAVALUE_CONTROLLED_CREATURES)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
