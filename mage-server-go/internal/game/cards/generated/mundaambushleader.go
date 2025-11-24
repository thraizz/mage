package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Munda Ambush Leader", NewMundaAmbushLeader)
}

// NewMundaAmbushLeader creates a Munda Ambush Leader
// {2}{R}{W} - CREATURE
// Haste
func NewMundaAmbushLeader(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Munda Ambush Leader")
	card.ManaCost = "{2}{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KOR", "ALLY"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(4, Integer.MAX_VALUE, filter, PutCards.TOP_ANY, Pu...)
	// card.AddAbility(ability1)
	return card, nil
}
