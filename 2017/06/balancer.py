found_inputs = {}
number_of_steps = 0


def balance_input(mbank):
    max_index = mbank.index(max(mbank))
    blocks = mbank[max_index]
    max_value = max(mbank)
    mbank[max_index] = 0
    tokens_left = max_value
    i = 1
    while tokens_left > 0:
        next_index = (max_index + i) % len(mbank)
        mbank[next_index] = mbank[next_index] + 1
        tokens_left = tokens_left - 1
        i = i + 1
    print ("mb", mbank)
    return mbank


with open("input", "r") as f:
    mbank = [int(x) for x in f.read().split("\t")]
    print("i", mbank)
    key = "".join([str(x) for x in mbank])

    while found_inputs.get(key, -1) != 1:
        found_inputs[key] = 1
        mbank = balance_input(mbank)
        key = "".join([str(x) for x in mbank])
        number_of_steps = number_of_steps + 1
    print("st", number_of_steps)
