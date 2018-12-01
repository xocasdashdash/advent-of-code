with open("input", "r") as f:
    input = f.read()
    lines = input.split("\n")
    diff = 0
    for line in lines:
        if not len(line):
            continue
        numbers = [int(x) for x in line.split("\t")]
        diff = diff + (max(numbers) - min(numbers))
    print("CS", diff)
