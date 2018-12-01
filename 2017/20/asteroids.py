data = open("input","r").read()
def parse(particle):
    return [list(map(int, p[3:-1].split(","))) for p in particle.split(", ")]

def step(d):
    d[1][0] += d[2][0]
    d[1][1] += d[2][1]
    d[1][2] += d[2][2]
    d[0][0] += d[1][0]
    d[0][1] += d[1][1]
    d[0][2] += d[1][2]

def part1(data):
    particles = [parse(d) for d in data.split('\n')]
    while True:
        for d in particles:
            step(d)
        m = sum([abs(e) for e in particles[0][0]])
        min_n = 0
        for i, d in enumerate(particles):
            if sum([abs(e) for e in d[0]]) < m:
                min_n = i
                m = sum([abs(e) for e in d[0]])
        print(min_n)

def part2(data):
    particles = [parse(d) for d in data.split('\n')]
    while True:
        positions = {}
        delete = []
        for i, d in enumerate(particles):
            step(d)
            if tuple(d[0]) in positions:
                delete += [i, positions[tuple(d[0])]]
            else:
                positions[tuple(d[0])] = i
        particles = [d for i, d in enumerate(particles) if i not in delete]
        print(len(particles))
#part1(data)
part2(data)
