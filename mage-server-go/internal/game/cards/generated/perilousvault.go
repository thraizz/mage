package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Perilous Vault", NewPerilousVault)
}

// NewPerilousVault creates a Perilous Vault
// {4} - ARTIFACT
func NewPerilousVault(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Perilous Vault")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{5}").
		AddTapCost().
		// TODO: ExileAllEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
