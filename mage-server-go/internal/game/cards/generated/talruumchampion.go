package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Talruum Champion", NewTalruumChampion)
}

// NewTalruumChampion creates a Talruum Champion
// {4}{R} - CREATURE
// FirstStrike, FirstStrike
func NewTalruumChampion(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Talruum Champion")
	card.ManaCost = "{4}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MINOTAUR"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability1)
	return card, nil
}
