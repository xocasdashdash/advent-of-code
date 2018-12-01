with open("input", "r") as f:
    input = f.read()
    lines = input.split("\n")
    valid = 0
    invalid = 0
    for line in lines:
        if len(line) == 0:
            continue
        words = [''.join(y) for y in [sorted(x) for x in line.split(" ")]]
        print(words)
        if len(words) == len(set(words)):
            valid = valid + 1
        else:
            invalid = invalid + 1
    print("valid", valid)
    print("invalid", invalid)
