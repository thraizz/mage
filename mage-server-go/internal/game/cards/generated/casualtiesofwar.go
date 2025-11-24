package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Casualties Of War", NewCasualtiesOfWar)
}

// NewCasualtiesOfWar creates a Casualties Of War
// {2}{B}{B}{G}{G} - SORCERY
func NewCasualtiesOfWar(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Casualties Of War")
	card.ManaCost = "{2}{B}{B}{G}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDestroyEffect()).
		AddEffect(abilities.NewDestroyEffect()).
		AddEffect(abilities.NewDestroyEffect()).
		AddEffect(abilities.NewDestroyEffect()).
		AddEffect(abilities.NewDestroyEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}