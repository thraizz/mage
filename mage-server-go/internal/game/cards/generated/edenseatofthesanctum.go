package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Eden Seat Of The Sanctum", NewEdenSeatOfTheSanctum)
}

// NewEdenSeatOfTheSanctum creates a Eden Seat Of The Sanctum
//  - LAND
func NewEdenSeatOfTheSanctum(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Eden Seat Of The Sanctum")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"TOWN"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{5}").
		AddTapCost().
		AddEffect(abilities.NewMillCardsControllerEffect(1)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}