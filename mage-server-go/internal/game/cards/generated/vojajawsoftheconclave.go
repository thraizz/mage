package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Voja Jaws Of The Conclave", NewVojaJawsOfTheConclave)
}

// NewVojaJawsOfTheConclave creates a Voja Jaws Of The Conclave
// {2}{R}{G}{W} - CREATURE
// Vigilance, Trample
func NewVojaJawsOfTheConclave(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Voja Jaws Of The Conclave")
	card.ManaCost = "{2}{R}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WOLF"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability1)
	return card, nil
}
