package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cursed Mirror", NewCursedMirror)
}

// NewCursedMirror creates a Cursed Mirror
// {2}{R} - ARTIFACT
// Haste
func NewCursedMirror(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cursed Mirror")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "R")
	card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - CopyPermanentEffect(new CursedMirrorCopyApplier())
	// card.AddAbility(ability2)
	return card, nil
}