package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Harald King Of Skemfar", NewHaraldKingOfSkemfar)
}

// NewHaraldKingOfSkemfar creates a Harald King Of Skemfar
// {1}{B}{G} - CREATURE
// Menace
func NewHaraldKingOfSkemfar(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Harald King Of Skemfar")
	card.ManaCost = "{1}{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "WARRIOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordMenace)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 5, 1, filter, PutCards.HAND, PutC...)
	// card.AddAbility(ability1)
	return card, nil
}
