with open("input", "r") as f:
    inp = [int(x) for x in f.read().split("\n")]
    curr_index = 0
    num_steps = 0
    while curr_index < len(inp):
        next_index = curr_index + inp[curr_index]
        num_steps = num_steps + 1
        inp[curr_index] = inp[curr_index] + 1
        curr_index = next_index
    print("ns", num_steps)
