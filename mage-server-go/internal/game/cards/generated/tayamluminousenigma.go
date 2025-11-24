package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tayam Luminous Enigma", NewTayamLuminousEnigma)
}

// NewTayamLuminousEnigma creates a Tayam Luminous Enigma
// {1}{W}{B}{G} - CREATURE
func NewTayamLuminousEnigma(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tayam Luminous Enigma")
	card.ManaCost = "{1}{W}{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"NIGHTMARE", "BEAST"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{3}").
		AddEffect(abilities.NewMillCardsControllerEffect(1)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
