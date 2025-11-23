package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nourishing Shoal", NewNourishingShoal)
}

// NewNourishingShoal creates a Nourishing Shoal
// {X}{G}{G} - INSTANT
func NewNourishingShoal(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nourishing Shoal")
	card.ManaCost = "{X}{G}{G}"
	card.Types = []string{"INSTANT"}
	card.Subtypes = []string{"ARCANE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainLifeEffect(ExileFromHandCostCardConvertedMana.instance)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
