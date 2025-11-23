package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Quandrix Cultivator", NewQuandrixCultivator)
}

// NewQuandrixCultivator creates a Quandrix Cultivator
// {1}{G}{G/U}{U} - CREATURE
func NewQuandrixCultivator(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Quandrix Cultivator")
	card.ManaCost = "{1}{G}{G/U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"TURTLE", "DRUID"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewSearchLibraryPutInPlayEffect(abilities.NewTargetRequirement(0, 1, abilities.NewAnyTargetFilter()), false)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
