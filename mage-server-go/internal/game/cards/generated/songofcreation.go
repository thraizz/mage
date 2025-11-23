package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Song Of Creation", NewSongOfCreation)
}

// NewSongOfCreation creates a Song Of Creation
// {1}{G}{U}{R} - ENCHANTMENT
func NewSongOfCreation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Song Of Creation")
	card.ManaCost = "{1}{G}{U}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(2)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DiscardHandControllerEffect()
	// card.AddAbility(ability1)
	return card, nil
}
