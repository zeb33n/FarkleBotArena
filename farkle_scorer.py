from collections import Counter

class DiceFondler:
    def __init__(self, roll, prev_score):
        self.roll = roll
        self.counts = Counter(roll)
        self.score = prev_score
        self.reroll = 1
        self.diceleft = 0
        self.calculate_score()

    def calculate_score(self):
        for value, count in self.counts.items():
            if count >= 4:
                self.score += 1000 * (count - 3)
            elif count == 3:
                if value == 1:
                    self.score += 300
                else:
                    self.score += value * 100
            elif value == 1:
                self.score += count * 100
            elif value == 5:
                self.score += count * 50
            else:
                self.reroll = 0
                self.diceleft += 1

        if list(self.roll) == [1, 2, 3, 4, 5, 6]:
            self.score, self.reroll, self.diceleft = 1500, 1, 0
            return

        if len(self.counts) == 3 and all(count == 2 for count in self.counts.values()):
            self.score, self.reroll, self.diceleft = 1500, 1, 0
            return

        if len(self.counts) == 2 and all(count == 3 for count in self.counts.values()):
            self.score, self.reroll, self.diceleft = 2500, 1, 0
            return

    def get_rollscore(self):
        return self.score - prev_score
    
    def get_totscore(self):
        return self.score

    def get_reroll(self):
        return self.reroll

    def get_diceleft(self):
        return self.diceleft

    def get_roll(self):
        return self.roll

#how to use
#not sure how rounds are being handled atm
roll = (1, 2, 3, 3, 4, 5)
prev_score = 0
round_results = DiceFondler(roll, prev_score)
print(f"player x on their n roll:{round_results.get_roll()} Scored: {round_results.get_rollscore()} points! with {round_results.get_diceleft()} dice left to roll, Total points this go: {round_results.get_totscore()} " )
