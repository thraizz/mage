package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Shaman En Kor", NewShamanEnKor)
}

// NewShamanEnKor creates a Shaman En Kor
// {1}{W} - CREATURE
func NewShamanEnKor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Shaman En Kor")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KOR", "CLERIC", "SHAMAN"}
	card.Power = "1"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - RedirectDamageFromSourceToTargetEffect()
	//
	// Costs:
	//   - AddManaCost("{0}")
	// card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - ShamanEnKorRedirectFromTargetEffect()
	// card.AddAbility(ability1)
	return card, nil
}
