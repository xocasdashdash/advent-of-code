from collections import Counter

with open("input", "r") as f:
    route = f.read().split(",")
directions = Counter(route)


def cancel(dir1, dir2):
    canceled = min(directions[dir1], directions[dir2])
    directions[dir1] -= canceled
    directions[dir2] -= canceled
    return canceled


done = False
max_length = 0


def calc_max(max_length):
    if max_length >= sum(directions.values()):
        print("change", max_length)
        max_length = sum(directions.values())


while not done:
    length = sum(directions.values())

    # Cancel opposing directions
    cancel('n', 's')
    calc_max(max_length)
    cancel('sw', 'ne')
    calc_max(max_length)
    cancel('se', 'nw')
    calc_max(max_length)
    directions['n'] += cancel('ne', 'nw')
    calc_max(max_length)
    directions['ne'] += cancel('se', 'n')
    calc_max(max_length)
    directions['se'] += cancel('ne', 's')
    calc_max(max_length)
    directions['s'] += cancel('se', 'sw')
    calc_max(max_length)
    directions['sw'] += cancel('s', 'nw')
    calc_max(max_length)
    directions['nw'] += cancel('n', 'sw')
    calc_max(max_length)

    new_length = sum(directions.values())
    if new_length >= max_length:
        max_length = new_length
    done = (length == new_length)

print(length)
print(max_length)
