with open("known_input", "r") as f:
    firewall_rules = [x for x in f.read().split("\n")]

firewall_steps = {}
max_level = 0
for f in firewall_rules:
    max_level = max_level if max_level >= int(
        "".join(f.split(":")[0])) else int("".join(f.split(":")[0]))
    firewall_steps[int("".join(f.split(":")[0]))] = int(
        "".join(f.split(":")[1:]))
print("fs", firewall_steps)
pos = 0
hits = []
for i in range(max_level + 1):
    if i not in firewall_steps:
        continue
    curr_step = firewall_steps[i]
    next_step = firewall_steps.get(i + 1, -1)
    if pos == 0:
        hit = (i, curr_step)
        print("hit!", hit)
        hits.append((i, curr_step))
print("hits", hits)
