package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lunatic Pandora", NewLunaticPandora)
}

// NewLunaticPandora creates a Lunatic Pandora
// {1} - ARTIFACT
func NewLunaticPandora(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lunatic Pandora")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{2}").
		AddTapCost().
		AddEffect(abilities.NewSurveilEffect(1)).
		Build()
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{6}").
		AddTapCost().
		AddSacrificeSourceCost().
		AddEffect(abilities.NewDestroyEffect()).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
