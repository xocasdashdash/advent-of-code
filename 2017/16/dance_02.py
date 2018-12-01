with open("input","r") as f:
    moves = f.read().strip().split(",")
dancers = (list(map(chr, range(97, 113))))
origDancers = (list(map(chr, range(97, 113))))
def switchIndexes (i, j, dancers):
    dancers[i], dancers[j]= dancers[j],dancers[i]
def switchElements (a,b, dancers):
    indexA = dancers.index(a)
    indexB = dancers.index(b)
    switchIndexes(indexA,indexB, dancers)
def dance(moves, dancers):
    for m in moves:
        move = m[0]
        if move == "s":
            f1 =  dancers[-1*int(m[1:]):]
            f2 = dancers[:len(dancers)-int(m[1:])]
            dancers = f1+f2
        elif move == "x":
            split = m[1:].split("/")
            switchIndexes(int(split[0]), int(split[1]), dancers)
        elif move == "p":
            switchElements(m[1], m[3], dancers)
    return dancers
rounds = 0
print("rest", 1000000001%30)
print("len", len(dancers))
print("r3st", 1000000000%(2*len(dancers)-2))
for i in range(10):
    dancers = dance(moves, dancers)
    rounds +=1
    if "".join(dancers) == "".join(origDancers):
        print("i",rounds)

print("dancers", "".join(dancers))
