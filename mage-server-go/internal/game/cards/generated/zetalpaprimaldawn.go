package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Zetalpa Primal Dawn", NewZetalpaPrimalDawn)
}

// NewZetalpaPrimalDawn creates a Zetalpa Primal Dawn
// {6}{W}{W} - CREATURE
// Flying, DoubleStrike, Vigilance, Trample, Indestructible
func NewZetalpaPrimalDawn(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Zetalpa Primal Dawn")
	card.ManaCost = "{6}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDER", "DINOSAUR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "8"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDoubleStrike)
	card.AddAbility(ability1)
	ability2 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability2)
	ability3 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability3)
	ability4 := abilities.NewKeywordAbility(card.ID, abilities.KeywordIndestructible)
	card.AddAbility(ability4)
	return card, nil
}
