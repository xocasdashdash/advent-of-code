cube_directions = {
    "n": {"x": 1, "y": -1, "z": 0},  # n
    "s": {"x": -1, "y": 1, "z": 0},  # s
    "sw": {"x": -1, "y": 0, "z": 1},  # sw
    "se": {"x": 0, "y": -1, "z": 1},  # se
    "ne": {"x": 1, "y": 0, "z": -1},  # ne
    "nw": {"x": 0, "y": 1, "z": -1},  # nw
}


def cube_add(cube, direction):
    return {
        "x": cube['x'] + cube_directions[direction]['x'],
        "y": cube['y'] + cube_directions[direction]['y'],
        "z": cube['z'] + cube_directions[direction]['z'],
    }


with open("input", "r") as f:
    route = f.read().split(",")
start_cube = {"x": 0, "y": 0, "z": 0}
max_value = 0

for r in route:
    start_cube = cube_add(start_cube, r)
    for k in "xyz":
        if start_cube[k] > max_value:
            print ("sc", start_cube)
            max_value = start_cube[k]

print(start_cube)
print(max_value)
