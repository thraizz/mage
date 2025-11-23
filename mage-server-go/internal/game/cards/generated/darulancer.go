package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Daru Lancer", NewDaruLancer)
}

// NewDaruLancer creates a Daru Lancer
// {4}{W}{W} - CREATURE
// FirstStrike
func NewDaruLancer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Daru Lancer")
	card.ManaCost = "{4}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "SOLDIER"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability0)
	return card, nil
}
