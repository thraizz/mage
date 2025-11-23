package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hraesvelgr Of The First Brood", NewHraesvelgrOfTheFirstBrood)
}

// NewHraesvelgrOfTheFirstBrood creates a Hraesvelgr Of The First Brood
// {4}{U} - CREATURE
// Flying, Vigilance
func NewHraesvelgrOfTheFirstBrood(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hraesvelgr Of The First Brood")
	card.ManaCost = "{4}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDER", "DRAGON"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability1)
	return card, nil
}
