package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kasla The Broken Halo", NewKaslaTheBrokenHalo)
}

// NewKaslaTheBrokenHalo creates a Kasla The Broken Halo
// {3}{U}{R}{W} - CREATURE
// Flying, Vigilance, Haste
func NewKaslaTheBrokenHalo(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kasla The Broken Halo")
	card.ManaCost = "{3}{U}{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ANGEL", "ALLY"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability1)
	ability2 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability2)
	return card, nil
}
