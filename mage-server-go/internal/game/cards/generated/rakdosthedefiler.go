package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rakdos The Defiler", NewRakdosTheDefiler)
}

// NewRakdosTheDefiler creates a Rakdos The Defiler
// {2}{B}{B}{R}{R} - CREATURE
// Flying, Trample
func NewRakdosTheDefiler(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rakdos The Defiler")
	card.ManaCost = "{2}{B}{B}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DEMON"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "7"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeControllerEffect(attackTriggerFilter, new HalfValue(new PermanentsO...)
	// card.AddAbility(ability2)
	return card, nil
}
