package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Korvold Fae Cursed King", NewKorvoldFaeCursedKing)
}

// NewKorvoldFaeCursedKing creates a Korvold Fae Cursed King
// {2}{B}{R}{G} - CREATURE
// Flying
func NewKorvoldFaeCursedKing(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Korvold Fae Cursed King")
	card.ManaCost = "{2}{B}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DRAGON", "NOBLE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeControllerEffect(filter, 1, "")
	// card.AddAbility(ability1)
	return card, nil
}
