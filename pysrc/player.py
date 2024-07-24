with open(
    __file__.rsplit("/", 2)[0] + "/" + __file__.rsplit("/", 1)[1].replace(".py", "")
) as p:
    while True:
        data = p.read()
        if data:
            print(data)
            break
