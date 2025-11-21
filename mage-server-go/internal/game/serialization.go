package game

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// SerializationChecksum computes a deterministic checksum of a game state snapshot
// Per Java SerializationTest and DeckHashTest: checksums ensure state integrity
// and guard against divergent game states across replays or network transmission
type SerializationChecksum struct {
	Hash      string // SHA-256 hash of deterministic serialization
	Timestamp string // ISO timestamp when checksum was computed
	Version   int    // Serialization version (for forward compatibility)
}

// ComputeChecksum generates a deterministic checksum of the game state snapshot
// The checksum is based on sorted, deterministic representation of all game state
// excluding non-deterministic fields like timestamps and random IDs
func (snapshot *gameStateSnapshot) ComputeChecksum() (*SerializationChecksum, error) {
	// Build deterministic representation
	// Sort all maps by keys, arrays by stable fields
	deterministicData := snapshot.buildDeterministicRepresentation()

	// Compute SHA-256 hash
	hash := sha256.New()
	if _, err := hash.Write([]byte(deterministicData)); err != nil {
		return nil, fmt.Errorf("failed to compute hash: %w", err)
	}

	return &SerializationChecksum{
		Hash:      hex.EncodeToString(hash.Sum(nil)),
		Timestamp: snapshot.Timestamp.Format("2006-01-02T15:04:05.000Z"),
		Version:   1,
	}, nil
}

// buildDeterministicRepresentation creates a canonical string representation
// of the game state that is independent of map iteration order or timestamps
func (snapshot *gameStateSnapshot) buildDeterministicRepresentation() string {
	var buf bytes.Buffer

	// Core game state (deterministic fields only)
	buf.WriteString(fmt.Sprintf("GAME:%s|%s|%d|%s|%s|%d\n",
		snapshot.GameID,
		snapshot.GameType,
		snapshot.State,
		snapshot.ActivePlayer,
		snapshot.PriorityPlayer,
		snapshot.TurnNumber,
	))

	// Players - sorted by ID
	playerIDs := make([]string, 0, len(snapshot.Players))
	for id := range snapshot.Players {
		playerIDs = append(playerIDs, id)
	}
	sort.Strings(playerIDs)

	for _, id := range playerIDs {
		player := snapshot.Players[id]
		buf.WriteString(fmt.Sprintf("PLAYER:%s|%s|%d|%d|%d|%t|%t|%d|%d\n",
			id,
			player.Name,
			player.Life,
			player.Poison,
			player.Energy,
			player.Passed,
			player.Lost,
			len(player.Library),
			len(player.Hand),
		))

		// Player graveyard - sorted by card ID for determinism
		graveyardIDs := make([]string, len(player.Graveyard))
		for i, card := range player.Graveyard {
			graveyardIDs[i] = card.ID
		}
		sort.Strings(graveyardIDs)
		for _, cardID := range graveyardIDs {
			buf.WriteString(fmt.Sprintf("  GRAVEYARD:%s\n", cardID))
		}
	}

	// Cards - sorted by ID
	cardIDs := make([]string, 0, len(snapshot.Cards))
	for id := range snapshot.Cards {
		cardIDs = append(cardIDs, id)
	}
	sort.Strings(cardIDs)

	for _, id := range cardIDs {
		card := snapshot.Cards[id]
		buf.WriteString(fmt.Sprintf("CARD:%s|%s|%s|%s|%d|%s|%s|%d|%t|%t\n",
			id,
			card.Name,
			card.OwnerID,
			card.ControllerID,
			card.Zone,
			card.Power,
			card.Toughness,
			card.Damage,
			card.Tapped,
			card.SummoningSickness,
		))

		// Card type (single string field)
		buf.WriteString(fmt.Sprintf("  TYPE:%s\n", card.Type))

		// Card subtypes - sorted
		subtypes := make([]string, len(card.SubTypes))
		copy(subtypes, card.SubTypes)
		sort.Strings(subtypes)
		for _, st := range subtypes {
			buf.WriteString(fmt.Sprintf("  SUBTYPE:%s\n", st))
		}

		// Card counters - sorted
		if card.Counters != nil {
			counterList := card.Counters.GetAll()
			counterNames := make([]string, 0, len(counterList))
			for _, counter := range counterList {
				counterNames = append(counterNames, counter.Name)
			}
			sort.Strings(counterNames)
			for _, name := range counterNames {
				for _, counter := range counterList {
					if counter.Name == name {
						buf.WriteString(fmt.Sprintf("  COUNTER:%s=%d\n", counter.Name, counter.Count))
						break
					}
				}
			}
		}

		// Card abilities - sorted by ability ID
		abilityIDs := make([]string, len(card.Abilities))
		for i, ability := range card.Abilities {
			abilityIDs[i] = ability.ID
		}
		sort.Strings(abilityIDs)
		for _, abilityID := range abilityIDs {
			buf.WriteString(fmt.Sprintf("  ABILITY:%s\n", abilityID))
		}
	}

	// Battlefield - sorted by card ID
	battlefieldIDs := make([]string, len(snapshot.Battlefield))
	for i, card := range snapshot.Battlefield {
		battlefieldIDs[i] = card.ID
	}
	sort.Strings(battlefieldIDs)
	buf.WriteString("BATTLEFIELD:")
	buf.WriteString(strings.Join(battlefieldIDs, ","))
	buf.WriteString("\n")

	// Exile - sorted by card ID
	exileIDs := make([]string, len(snapshot.Exile))
	for i, card := range snapshot.Exile {
		exileIDs[i] = card.ID
	}
	sort.Strings(exileIDs)
	buf.WriteString("EXILE:")
	buf.WriteString(strings.Join(exileIDs, ","))
	buf.WriteString("\n")

	// Command - sorted by card ID
	commandIDs := make([]string, len(snapshot.Command))
	for i, card := range snapshot.Command {
		commandIDs[i] = card.ID
	}
	sort.Strings(commandIDs)
	buf.WriteString("COMMAND:")
	buf.WriteString(strings.Join(commandIDs, ","))
	buf.WriteString("\n")

	// Stack - order matters for stack (LIFO), so don't sort
	buf.WriteString("STACK:\n")
	for i, item := range snapshot.StackItems {
		buf.WriteString(fmt.Sprintf("  %d:%s|%s\n", i, item.ID, item.Description))
	}

	// Player order - order matters
	buf.WriteString("PLAYER_ORDER:")
	buf.WriteString(strings.Join(snapshot.PlayerOrder, ","))
	buf.WriteString("\n")

	return buf.String()
}

// VerifyChecksum verifies that a snapshot's stored checksum matches its computed checksum
// Returns true if checksums match, false otherwise
func (snapshot *gameStateSnapshot) VerifyChecksum(expected *SerializationChecksum) (bool, error) {
	computed, err := snapshot.ComputeChecksum()
	if err != nil {
		return false, fmt.Errorf("failed to compute checksum: %w", err)
	}

	return computed.Hash == expected.Hash, nil
}

// SerializeToBytes serializes a game state snapshot to bytes using gob encoding
// This is the standard serialization used for replay files and network transmission
func (snapshot *gameStateSnapshot) SerializeToBytes() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(snapshot); err != nil {
		return nil, fmt.Errorf("failed to encode snapshot: %w", err)
	}

	return buf.Bytes(), nil
}

// DeserializeFromBytes deserializes a game state snapshot from bytes using gob encoding
func DeserializeFromBytes(data []byte) (*gameStateSnapshot, error) {
	var snapshot gameStateSnapshot
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)

	if err := decoder.Decode(&snapshot); err != nil {
		return nil, fmt.Errorf("failed to decode snapshot: %w", err)
	}

	return &snapshot, nil
}

// ValidateSerializationRoundtrip validates that a snapshot can be serialized
// and deserialized without data loss by comparing checksums
func ValidateSerializationRoundtrip(snapshot *gameStateSnapshot) error {
	// Compute checksum of original
	originalChecksum, err := snapshot.ComputeChecksum()
	if err != nil {
		return fmt.Errorf("failed to compute original checksum: %w", err)
	}

	// Serialize
	data, err := snapshot.SerializeToBytes()
	if err != nil {
		return fmt.Errorf("failed to serialize: %w", err)
	}

	// Deserialize
	deserialized, err := DeserializeFromBytes(data)
	if err != nil {
		return fmt.Errorf("failed to deserialize: %w", err)
	}

	// Compute checksum of deserialized
	deserializedChecksum, err := deserialized.ComputeChecksum()
	if err != nil {
		return fmt.Errorf("failed to compute deserialized checksum: %w", err)
	}

	// Compare checksums
	if originalChecksum.Hash != deserializedChecksum.Hash {
		return fmt.Errorf("checksum mismatch: original=%s, deserialized=%s",
			originalChecksum.Hash, deserializedChecksum.Hash)
	}

	return nil
}
