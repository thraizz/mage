package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Eyeblights Ending", NewEyeblightsEnding)
}

// NewEyeblightsEnding creates a Eyeblights Ending
// {2}{B} - KINDRED INSTANT
func NewEyeblightsEnding(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Eyeblights Ending")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"KINDRED", "INSTANT"}
	card.Subtypes = []string{"ELF"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDestroyEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
