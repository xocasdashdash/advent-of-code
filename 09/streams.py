with open("input", "r") as fich:
    bad_input = fich.read()

good_input = []
noise_on = False
i = 0
garbage = 0
garbage_counter_on = False
score = 0
depth = 0
while i < len(bad_input):
    garbage_counter_on = noise_on
    curr_char = bad_input[i]
    if curr_char == '!':
        i = i + 2
        continue
    if not noise_on and curr_char not in "<>!,":
        good_input.append(curr_char)
    if curr_char == '<':
        noise_on = True
    elif curr_char == '>':
        noise_on = False
    if garbage_counter_on and noise_on:
        garbage = garbage + 1
    i = i + 1
    if not noise_on:
        if curr_char == "{":
            depth = depth + 1
            score = depth + score
        elif curr_char == "}":
            depth = depth - 1
clean_river = "{msg}".format(msg="".join(good_input))
print("score", score)
print("garbage", garbage)
