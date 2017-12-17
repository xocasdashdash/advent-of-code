

from functools import reduce

def knothash(lens):
    lens = [ord(x) for x in lens.rstrip()]
    lens.extend([17,31,73,47,23])
    nums = [x for x in range(0,256)]
    pos = 0
    skip = 0
    for _ in range(64):
        for l in lens:
            to_reverse = []
            for x in range(l):
                n = (pos + x) % 256
                to_reverse.append(nums[n])
            to_reverse.reverse()
            for x in range(l):
                n = (pos + x) % 256
                nums[n] = to_reverse[x]
            pos += l + skip
            pos = pos % 256
            skip += 1
    dense = []
    for x in range(0,16):
        subslice = nums[16*x:16*x+16]
        dense.append('%02x'%reduce((lambda x,y: x ^ y),subslice))
    return ''.join(dense)

def solve(key_string):
    count = 0
    unseen = []
    for i in range(128):
        hash = knothash(key_string + "-" + str(i))
        bin_hash = bin(int(hash, 16))[2:].zfill(128)
        unseen += [(i, j) for j, d in enumerate(bin_hash) if d == '1']
    print("Part 1: " + str(len(unseen)))
    while unseen:
        queued = [unseen[0]]
        while queued:
            (x, y) = queued.pop()
            if (x, y) in unseen:
                unseen.remove((x, y))
                queued += [(x - 1, y), (x+ 1, y), (x, y+1), (x, y-1)]
        count += 1
    print("Part 2: " + str(count))
solve("hxtvlmkl")
