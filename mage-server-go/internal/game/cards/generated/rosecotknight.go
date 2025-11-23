package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Rosecot Knight", NewRosecotKnight)
}

// NewRosecotKnight creates a Rosecot Knight
// {4}{W} - CREATURE
// Vigilance
func NewRosecotKnight(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rosecot Knight")
	card.ManaCost = "{4}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "KNIGHT"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 6, 1, filter, PutCards.HAND, PutC...)
	// card.AddAbility(ability1)
	return card, nil
}
