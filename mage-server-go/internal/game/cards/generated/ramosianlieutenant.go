package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ramosian Lieutenant", NewRamosianLieutenant)
}

// NewRamosianLieutenant creates a Ramosian Lieutenant
// {1}{W} - CREATURE
func NewRamosianLieutenant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ramosian Lieutenant")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "REBEL"}
	card.Power = "1"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddManaCost("{4}").
		// TODO: SearchLibraryPutInPlayEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
