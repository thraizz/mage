package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Engineered Explosives", NewEngineeredExplosives)
}

// NewEngineeredExplosives creates a Engineered Explosives
// {X} - ARTIFACT
func NewEngineeredExplosives(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Engineered Explosives")
	card.ManaCost = "{X}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddSacrificeSourceCost().
		AddEffect(abilities.NewDestroyEffect(filter)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
