key = "hxtvlmkl"
def calc_hash(key):
    lengths = [ord(x) for x in key] + [17, 31, 73, 47, 23]
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
    return "".join([x[2:] for x in map(bin, result)])
hashes = []
count = 0
for i in range(128):
    new_hash = calc_hash(key + "-" + str(i))
    count += len([x for x in new_hash if x == "1"])
    hashes.append(new_hash)
print(count)
for h in hashes:
    print("{}\n".format(h))
import sys
sys.setrecursionlimit(1500)

seen = set()
n = 0
def dfs(i, j):
    if ((i, j)) in seen:
        return
    if not hashes[i][j]:
        return
    seen.add((i, j))
    if i > 0:
        dfs(i-1, j)
    if j > 0:
        dfs(i, j-1)
    if i < 127:
        dfs(i+1, j)
    if j < 127:
        dfs(i, j+1)

for i in range(128):
    for j in range(128):
        if (i,j) in seen:
            continue
        if not hashes[i][j]:
            continue
        n += 1
        dfs(i, j)
