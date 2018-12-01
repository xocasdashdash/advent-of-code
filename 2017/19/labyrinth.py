lines = open("input").read().splitlines()
road = {x+1j*y: v for y, line in enumerate(lines) for x, v in enumerate(line) if v.strip()}
direction, pos, path = 1j, min(road, key=lambda v: v.imag), []
while pos in road:
    if road[pos] == '+':
        direction = next(d for d in [direction*1j, direction*-1j]
                         if pos+d in road and d != path[-1]-pos)
    path += [pos]
    pos += direction

print(''.join(c for c in map(road.get, path) if c.isalpha()))
print(len(path))
