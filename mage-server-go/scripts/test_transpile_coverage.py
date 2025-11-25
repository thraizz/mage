#!/usr/bin/env python3
"""
Test script to verify that the transpiler finds all Java card files.

This test:
1. Scans the Java cards directory to find all card files
2. Runs the transpiler on all cards
3. Verifies that all cards are processed (either successfully or with errors)
4. Reports any cards that were completely missed
"""

import os
import sys
import json
from pathlib import Path
from typing import Set, Dict, List
from collections import defaultdict

# Add scripts directory to path so we can import transpile_cards
sys.path.insert(0, str(Path(__file__).parent))

from transpile_cards import transpile_card, TranspileResult


def find_all_java_cards(java_cards_dir: str) -> Set[str]:
    """Find all Java card files and return their class names."""
    java_files = set()
    java_path = Path(java_cards_dir)
    
    if not java_path.exists():
        print(f"Error: Java cards directory not found: {java_cards_dir}")
        return java_files
    
    # Walk through all subdirectories
    for root, dirs, files in os.walk(java_path):
        for file in files:
            if file.endswith('.java'):
                # Extract class name from filename (remove .java extension)
                class_name = file[:-5]  # Remove .java
                java_files.add(class_name)
    
    return java_files


def get_processed_cards(results: List[TranspileResult]) -> Set[str]:
    """Extract card class names from transpilation results."""
    processed = set()
    for result in results:
        # Extract class name from the Java file path
        # The java_file field contains the full path
        if result.java_file:
            # Get the filename without extension
            class_name = Path(result.java_file).stem
            # Normalize to lowercase for comparison
            processed.add(class_name.lower())
    return processed


def run_transpile_test(java_cards_dir: str, output_dir: str, limit: int = 0, sample: int = 0) -> Dict:
    """Run the transpiler on all cards and collect results."""
    print("=" * 80)
    print("Transpiler Coverage Test")
    print("=" * 80)
    print()
    
    # Step 1: Find all Java card files
    print("Step 1: Scanning Java cards directory...")
    all_java_cards = find_all_java_cards(java_cards_dir)
    print(f"Found {len(all_java_cards)} Java card files")
    print()
    
    if not all_java_cards:
        print("ERROR: No Java card files found!")
        return {
            "success": False,
            "total_java_cards": 0,
            "processed_cards": 0,
            "missing_cards": [],
            "error": "No Java card files found"
        }
    
    # Step 2: Run transpiler on all cards
    print("Step 2: Running transpiler on all cards...")
    print("This may take a while...")
    print()
    
    results = []
    java_files = []
    
    # Collect all Java file paths
    java_path = Path(java_cards_dir)
    for root, dirs, files in os.walk(java_path):
        for file in sorted(files):
            if file.endswith('.java'):
                java_files.append(os.path.join(root, file))
    
    # Apply limit or sample if specified
    if limit and limit > 0:
        java_files = java_files[:limit]
        print(f"  Limited to first {limit} cards")
    elif sample and sample > 0:
        import random
        if sample < len(java_files):
            java_files = random.sample(java_files, sample)
            print(f"  Testing random sample of {sample} cards")
        else:
            print(f"  Sample size ({sample}) >= total cards, testing all")
    
    # Process each file
    for i, java_file in enumerate(java_files, 1):
        if i % 100 == 0:
            print(f"  Progress: {i}/{len(java_files)} cards processed...")
        
        result = transpile_card(java_file, output_dir)
        results.append(result)
    
    print(f"  Completed: {len(java_files)} cards processed")
    print()
    
    # Step 3: Analyze results
    print("Step 3: Analyzing results...")
    
    # Get all processed card names (normalized to lowercase)
    processed_cards = get_processed_cards(results)
    
    # Get the Java card names for the files we actually processed
    processed_java_files = set()
    for java_file in java_files:
        # Extract class name from file path
        class_name = Path(java_file).stem
        processed_java_files.add(class_name.lower())
    
    # Get all Java card names (normalized to lowercase) - only for comparison if we processed all
    java_card_names = {name.lower() for name in all_java_cards}
    
    # Find missing cards:
    # - If we processed all cards: cards that exist in Java but weren't processed
    # - If we processed a subset: cards in the subset that weren't processed
    if (not limit or limit == 0) and (not sample or sample == 0):
        # Processed all cards - compare against all Java cards
        missing_cards = java_card_names - processed_cards
    else:
        # Processed a subset - only compare against the subset we tried to process
        missing_cards = processed_java_files - processed_cards
    
    # Find cards that were processed but don't exist in Java (shouldn't happen)
    extra_cards = processed_cards - java_card_names
    
    # Categorize results
    successful = [r for r in results if r.success]
    failed = [r for r in results if not r.success]
    
    # Group failures by error type
    failures_by_type = defaultdict(list)
    for result in failed:
        error_type = result.error_type or "Unknown"
        failures_by_type[error_type].append(result.card_name)
    
    # Print summary
    print("=" * 80)
    print("Test Results Summary")
    print("=" * 80)
    if limit and limit > 0:
        print(f"Test mode:                Limited to first {limit} cards")
    elif sample and sample > 0:
        print(f"Test mode:                Random sample of {sample} cards")
    else:
        print(f"Test mode:                All cards")
    print(f"Total Java card files:     {len(all_java_cards)}")
    print(f"Cards attempted:           {len(java_files)}")
    print(f"Cards processed:           {len(processed_cards)}")
    print(f"  ✓ Successful:            {len(successful)}")
    print(f"  ✗ Failed:                {len(failed)}")
    print(f"Missing cards:             {len(missing_cards)}")
    if (not limit or limit == 0) and (not sample or sample == 0):
        print(f"Extra cards (unexpected):  {len(extra_cards)}")
    print()
    
    # Report missing cards
    if missing_cards:
        print("=" * 80)
        print("⚠️  MISSING CARDS (Found in Java but not processed by transpiler)")
        print("=" * 80)
        for card in sorted(missing_cards):
            print(f"  - {card}")
        print()
    else:
        print("✓ All Java cards were processed by the transpiler!")
        print()
    
    # Report extra cards (shouldn't happen, but good to know)
    if extra_cards:
        print("=" * 80)
        print("⚠️  EXTRA CARDS (Processed but not found in Java directory)")
        print("=" * 80)
        for card in sorted(extra_cards):
            print(f"  - {card}")
        print()
    
    # Report failure breakdown
    if failures_by_type:
        print("=" * 80)
        print("Failure Breakdown by Error Type")
        print("=" * 80)
        for error_type, cards in sorted(failures_by_type.items(), key=lambda x: -len(x[1])):
            print(f"{error_type}: {len(cards)} cards")
            if len(cards) <= 10:
                for card in sorted(cards):
                    print(f"  - {card}")
            else:
                print(f"  (showing first 10 of {len(cards)})")
                for card in sorted(cards)[:10]:
                    print(f"  - {card}")
        print()
    
    # Final verdict
    print("=" * 80)
    if missing_cards:
        print("❌ TEST FAILED: Some cards were not processed")
        print(f"   Missing {len(missing_cards)} card(s)")
        success = False
    else:
        print("✓ TEST PASSED: All Java cards were found and processed")
        success = True
    print("=" * 80)
    print()
    
    return {
        "success": success,
        "total_java_cards": len(all_java_cards),
        "processed_cards": len(processed_cards),
        "successful": len(successful),
        "failed": len(failed),
        "missing_cards": sorted(missing_cards),
        "extra_cards": sorted(extra_cards),
        "failures_by_type": {k: sorted(v) for k, v in failures_by_type.items()}
    }


def main():
    """Main test function."""
    import argparse
    
    parser = argparse.ArgumentParser(
        description='Test that transpiler finds all Java card files'
    )
    parser.add_argument(
        '--input',
        default='/Users/aron/dev/opensource/mage/Mage.Sets/src/mage/cards',
        help='Input directory with Java card files'
    )
    parser.add_argument(
        '--output',
        default='internal/game/cards/generated',
        help='Output directory for Go files (will be cleaned)'
    )
    parser.add_argument(
        '--json',
        help='Save results to JSON file'
    )
    parser.add_argument(
        '--limit',
        type=int,
        default=0,
        help='Limit number of cards to test (0 for all cards, useful for quick testing)'
    )
    parser.add_argument(
        '--sample',
        type=int,
        help='Test a random sample of N cards (faster for quick checks)'
    )
    
    args = parser.parse_args()
    
    # Clean output directory to ensure fresh test
    output_path = Path(args.output)
    if output_path.exists():
        import shutil
        print(f"Cleaning output directory: {output_path}")
        shutil.rmtree(output_path)
    output_path.mkdir(parents=True, exist_ok=True)
    
    # Run the test
    results = run_transpile_test(args.input, args.output, args.limit, args.sample)
    
    # Save results to JSON if requested
    if args.json:
        with open(args.json, 'w') as f:
            json.dump(results, f, indent=2)
        print(f"Results saved to: {args.json}")
        print()
    
    # Exit with appropriate code
    sys.exit(0 if results["success"] else 1)


if __name__ == '__main__':
    main()

