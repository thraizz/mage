package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Sandals Of Abdallah", NewSandalsOfAbdallah)
}

// NewSandalsOfAbdallah creates a Sandals Of Abdallah
// {4} - ARTIFACT
func NewSandalsOfAbdallah(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sandals Of Abdallah")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{2}").
		AddTapCost().
		AddEffect(abilities.NewGrantAbilityEffect("IslandwalkAbility")).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
