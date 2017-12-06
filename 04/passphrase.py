with open("input", "r") as f:
    input = f.read()
    lines = input.split("\n")
    valid = 0
    invalid = 0
    for line in lines:
        words = line.split(" ")
        if len(line) == 0:
            continue
        if len(words) == len(set(words)):
            valid = valid + 1
        else:
            invalid = invalid + 1
    print("valid", valid)
    print("invalid", invalid)
