from collections import Counter


def calculate_score(roll: tuple[int]) -> tuple[int, int]:
    score = 0
    counts = Counter(roll)
    num_dice = len(roll)

    if list(roll) == [1, 2, 3, 4, 5, 6]:
        return 1500, 0

    if len(counts) == 3 and all(count == 2 for count in counts.values()):
        return 1500, 0

    if len(counts) == 2 and all(count == 3 for count in counts.values()):
        return 2500, 0

    for value, count in counts.items():
        if count >= 4:
            score += 1000 * (count - 3)
            num_dice -= count
        elif count == 3:
            if value == 1:
                score += 300
            else:
                score += value * 100
            num_dice -= count
        elif value == 1:
            score += count * 100
            num_dice -= count
        elif value == 5:
            score += count * 50
            num_dice -= count

    return score, num_dice
