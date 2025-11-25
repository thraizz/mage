package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Flame Of Anor", NewFlameOfAnor)
}

// NewFlameOfAnor creates a Flame Of Anor
// {1}{U}{R} - INSTANT
func NewFlameOfAnor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Flame Of Anor")
	card.ManaCost = "{1}{U}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(5)).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		AddEffect(abilities.NewDestroyEffect()).
		AddTarget(abilities.NewPlayerTargetFilter()).
		AddTarget(abilities.NewArtifactTargetFilter()).
		AddTarget(abilities.NewCreatureTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
