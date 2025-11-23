package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Pyrewood Gearhulk", NewPyrewoodGearhulk)
}

// NewPyrewoodGearhulk creates a Pyrewood Gearhulk
// {2}{R}{R}{G}{G} - ARTIFACT CREATURE
// Vigilance, Menace
func NewPyrewoodGearhulk(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Pyrewood Gearhulk")
	card.ManaCost = "{2}{R}{R}{G}{G}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"CONSTRUCT"}
	card.Power = "7"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordMenace)
	card.AddAbility(ability1)
	return card, nil
}
