package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Steward Of Valeron", NewStewardOfValeron)
}

// NewStewardOfValeron creates a Steward Of Valeron
// {G}{W} - CREATURE
// Vigilance
func NewStewardOfValeron(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Steward Of Valeron")
	card.ManaCost = "{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "DRUID", "KNIGHT"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "G")
	card.AddAbility(ability1)
	return card, nil
}
