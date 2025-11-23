package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Galea Kindler Of Hope", NewGaleaKindlerOfHope)
}

// NewGaleaKindlerOfHope creates a Galea Kindler Of Hope
// {1}{G}{W}{U} - CREATURE
// Vigilance
func NewGaleaKindlerOfHope(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Galea Kindler Of Hope")
	card.ManaCost = "{1}{G}{W}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "KNIGHT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	return card, nil
}
