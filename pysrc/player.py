player_num = __file__.replace(".py", "")[-1]
with open(__file__.rsplit("/", 2)[0] + f"/pipes/splayer{player_num}") as p:
    while True:
        data = p.read()
        if data:
            print(data)
            break
