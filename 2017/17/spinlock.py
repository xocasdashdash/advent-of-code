step_size=367
p1=0
np =0
for t in range(1,50000000+1):
    np = ( np + step_size) % t +1
    if np == 1:
        p1 = t
print("p1",p1)
