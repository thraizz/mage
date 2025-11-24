package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fomori Vault", NewFomoriVault)
}

// NewFomoriVault creates a Fomori Vault
//  - LAND
func NewFomoriVault(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fomori Vault")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 ArtifactYouControlCount.instance,...)
	//
	// Costs:
	//   - AddManaCost("{3}")
	//   - AddTapCost()
	// card.AddAbility(ability1)
	return card, nil
}