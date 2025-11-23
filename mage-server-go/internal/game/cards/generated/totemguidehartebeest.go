package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Totem Guide Hartebeest", NewTotemGuideHartebeest)
}

// NewTotemGuideHartebeest creates a Totem Guide Hartebeest
// {4}{W} - CREATURE
func NewTotemGuideHartebeest(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Totem Guide Hartebeest")
	card.ManaCost = "{4}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ANTELOPE"}
	card.Power = "2"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewSearchLibraryPutInHandEffect(abilities.NewTargetRequirement(0, 1, abilities.NewAnyTargetFilter()), true)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
