def find_divisor(numbers):
    for index, number in enumerate(numbers):
        print("len", len(numbers[index + 1:]))
        for divider in reversed(numbers[index + 1:]):
            if number % divider == 0:
                print("found {} and {}. Rest: {}".format(
                    number, divider, number % divider))
                return int(number / divider)
    return 0


with open("input", "r") as f:
    input = f.read()
    lines = input.split("\n")
    sum = 0
    for line in lines:
        if not len(line):
            continue
        numbers = sorted([int(x) for x in line.split("\t")], reverse=True)
        sum = sum + find_divisor(numbers)
    print("CS", sum)
