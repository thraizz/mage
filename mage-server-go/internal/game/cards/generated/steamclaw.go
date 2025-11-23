package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Steamclaw", NewSteamclaw)
}

// NewSteamclaw creates a Steamclaw
// {2} - ARTIFACT
func NewSteamclaw(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Steamclaw")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddManaCost("{3}").
		AddEffect(abilities.NewExileTargetEffect()).
		Build()
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{1}").
		AddSacrificeSourceCost().
		AddEffect(abilities.NewExileTargetEffect()).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
