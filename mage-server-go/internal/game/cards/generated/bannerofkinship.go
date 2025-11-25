package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Banner Of Kinship", NewBannerOfKinship)
}

// NewBannerOfKinship creates a Banner Of Kinship
// {5} - ARTIFACT
func NewBannerOfKinship(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Banner Of Kinship")
	card.ManaCost = "{5}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: AsEntersBattlefieldAbility
	//   - Effect: ChooseCreatureTypeEffect()
	// card.AddAbility(ability0)
	return card, nil
}
