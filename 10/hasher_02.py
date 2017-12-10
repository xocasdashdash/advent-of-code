with open('input', 'r') as f:
    lengths = [ord(x) for x in f.read()] + [17, 31, 73, 47, 23]

inp = [x for x in range(0, 256)]
curr_position = 0
skip_size = 0
for i in range(64):
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
result = []
for i in range(16):
    start_index = i * 16
    end_index = start_index + 16
    sub_list = inp[start_index:end_index]
    res = 0
    for i in range(0, 16):
        res ^= sub_list[i]
    result.append(res)
print("result", "".join([x[2:] for x in map(hex, result)]))
