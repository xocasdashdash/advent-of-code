
DIVIDER = 2147483647
def next_number_gen(start,factor):
    while True:
        nv = (start * factor) % DIVIDER
        yield nv
        start = nv
def next_number_gen_thatdivs(start,factor,div):
    while True:
        nv = next(next_number_gen(start,factor))
        if nv % div == 0:
            yield nv
        start = nv


def generate(start_a, start_b, n):
    side_a_gen = next_number_gen_thatdivs(start_a,16807,4)
    side_b_gen = next_number_gen_thatdivs(start_b,48271,8)

    count = 0
    for i in range(0,n):
        a = hex(next(side_a_gen))[-4:]
        b = hex(next(side_b_gen))[-4:]
        if a == b:
            count+=1
    print("count", count)
#generate(65, 8921,1057)
#generate(65, 8921,5000000)
generate(634, 301,5000000)
