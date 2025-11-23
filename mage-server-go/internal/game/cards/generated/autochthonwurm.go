package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Autochthon Wurm", NewAutochthonWurm)
}

// NewAutochthonWurm creates a Autochthon Wurm
// {10}{G}{G}{G}{W}{W} - CREATURE
// Trample
func NewAutochthonWurm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Autochthon Wurm")
	card.ManaCost = "{10}{G}{G}{G}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WURM"}
	card.Power = "9"
	card.Toughness = "14"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
