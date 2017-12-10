with open('input', 'r') as f:
    lengths = list(map(int, map(str.strip, f.read().split(','))))
inp = [x for x in range(0, 256)]
curr_position = 0
skip_size = 0
for length in lengths:
    start_sublist = []
    if curr_position + length >= len(inp):
        start_sublist = inp[0:curr_position + length - len(inp)]
    reversed_sublist = list(
        reversed(inp[curr_position:curr_position + length] + start_sublist))
    start_position = curr_position
    for i in range(start_position, start_position + len(reversed_sublist)):
        inp[i % len(inp)] = reversed_sublist[i - start_position]
    curr_position = (curr_position + length + skip_size) % len(inp)
    skip_size += 1
print("result", inp[0] * inp[1])
