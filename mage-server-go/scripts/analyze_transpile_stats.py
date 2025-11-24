#!/usr/bin/env python3
"""
Analyze transpile_stats.json file
Provides various views and queries of transpilation statistics
"""

import json
import sys
import argparse
from pathlib import Path
from collections import Counter


def load_stats(stats_file: str) -> dict:
    """Load statistics from JSON file"""
    stats_path = Path(stats_file)
    if not stats_path.exists():
        print(f"Error: Stats file not found: {stats_file}")
        sys.exit(1)

    with open(stats_path) as f:
        return json.load(f)


def print_summary(stats: dict):
    """Print summary statistics"""
    summary = stats["summary"]
    print("=== Summary ===")
    print(f"Total cards:          {summary['total']}")
    print(f"✓ Successful:         {summary['successful']} ({summary['success_rate']})")
    print(f"✗ Failed:             {summary['failed']}")
    print(f"⚠ Has TODO:           {summary['has_todo']} ({summary['todo_rate']})")
    print(f"✅ Fully implemented: {summary['fully_implemented']} ({summary['complete_rate']})")


def print_errors(stats: dict):
    """Print error breakdown"""
    error_categories = stats.get("error_categories", {})

    if not error_categories:
        print("No errors found!")
        return

    print("\n=== Error Breakdown ===")
    for error_type, data in sorted(error_categories.items(), key=lambda x: -x[1]["count"]):
        print(f"\n{error_type}: {data['count']} cards")
        print(f"Cards: {', '.join(data['cards'][:10])}")
        if len(data['cards']) > 10:
            print(f"  ... and {len(data['cards']) - 10} more")


def print_todos(stats: dict, limit: int = 20):
    """Print TODO analysis"""
    todo_analysis = stats.get("todo_analysis", {})

    if todo_analysis.get("total", 0) == 0:
        print("No TODOs found!")
        return

    print(f"\n=== TODO Analysis ===")
    print(f"Cards with TODOs: {todo_analysis['total']}")

    cards_with_todos = todo_analysis.get("cards_with_todos", [])
    print(f"\nCards with TODOs (showing first {min(limit, len(cards_with_todos))}):")
    for card_name in cards_with_todos[:limit]:
        print(f"  - {card_name}")

    if len(cards_with_todos) > limit:
        print(f"  ... and {len(cards_with_todos) - limit} more")

    # Show sample details
    sample_details = todo_analysis.get("sample_details", {})
    if sample_details:
        print(f"\nSample TODO messages:")
        for card_name, details in list(sample_details.items())[:5]:
            print(f"\n{card_name} ({details['count']} TODO(s)):")
            for msg in details.get("messages", [])[:3]:
                print(f"  - {msg}")


def print_failed_cards(stats: dict):
    """Print detailed failed card information"""
    failed_cards = stats.get("failed_cards", [])

    if not failed_cards:
        print("No failed cards!")
        return

    print(f"\n=== Failed Cards ({len(failed_cards)}) ===")
    for card in failed_cards[:20]:
        print(f"\n{card['card_name']}:")
        print(f"  File: {card['java_file']}")
        print(f"  Error: {card['error_type']}")
        print(f"  Message: {card['error_message']}")

    if len(failed_cards) > 20:
        print(f"\n... and {len(failed_cards) - 20} more (see JSON file for full list)")


def export_csv(stats: dict, output_file: str):
    """Export failed cards to CSV"""
    failed_cards = stats.get("failed_cards", [])

    if not failed_cards:
        print("No failed cards to export!")
        return

    with open(output_file, 'w') as f:
        f.write("Card Name,Error Type,Error Message,Java File\n")
        for card in failed_cards:
            card_name = card['card_name'].replace(',', ';')
            error_type = card['error_type'].replace(',', ';')
            error_msg = card['error_message'].replace(',', ';').replace('\n', ' ')
            java_file = card['java_file'].replace(',', ';')
            f.write(f"{card_name},{error_type},{error_msg},{java_file}\n")

    print(f"Exported {len(failed_cards)} failed cards to {output_file}")


def main():
    parser = argparse.ArgumentParser(description='Analyze transpile statistics')
    parser.add_argument('--stats', default='transpile_stats.json',
                       help='Statistics file to analyze')
    parser.add_argument('--summary', action='store_true',
                       help='Show summary only')
    parser.add_argument('--errors', action='store_true',
                       help='Show error breakdown')
    parser.add_argument('--todos', action='store_true',
                       help='Show TODO analysis')
    parser.add_argument('--failed', action='store_true',
                       help='Show detailed failed cards')
    parser.add_argument('--export-csv', metavar='FILE',
                       help='Export failed cards to CSV file')
    parser.add_argument('--limit', type=int, default=20,
                       help='Limit number of items to show (default: 20)')

    args = parser.parse_args()

    # Load stats
    stats = load_stats(args.stats)

    # If no specific option, show all
    show_all = not (args.summary or args.errors or args.todos or args.failed or args.export_csv)

    if show_all or args.summary:
        print_summary(stats)

    if show_all or args.errors:
        print_errors(stats)

    if show_all or args.todos:
        print_todos(stats, args.limit)

    if show_all or args.failed:
        print_failed_cards(stats)

    if args.export_csv:
        export_csv(stats, args.export_csv)


if __name__ == '__main__':
    main()
