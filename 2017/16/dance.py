with open("input","r") as f:
    moves = f.read().strip().split(",")
dancers = (list(map(chr, range(97, 113))))

def switchIndexes (i, j):
    dancers[i], dancers[j]= dancers[j],dancers[i]
def switchElements (a,b):
    indexA = dancers.index(a)
    indexB = dancers.index(b)
    switchIndexes(indexA,indexB)
for m in moves:
    move = m[0]
    if move == "s":
        f1 =  dancers[-1*int(m[1:]):]
        f2 = dancers[:len(dancers)-int(m[1:])]
        print("f1",f1)
        print("f2",f2)
        dancers = f1+f2
    elif move == "x":
        split = m[1:].split("/")
        switchIndexes(int(split[0]), int(split[1]))
    elif move == "p":
        switchElements(m[1], m[3])
    #print("move",m)
    print("dancers", "".join(dancers))
    #input()

