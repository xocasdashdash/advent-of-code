import math
input = 361527
size = math.floor(math.sqrt(input)) + 2
print("size", size)
start = 1
size_i = 1
number_of_increments = 1
while True:
    size_i = size_i + 2
    start = int(math.pow((size_i), 2))
    if start >= input:
        break
    number_of_increments = number_of_increments + 1
print(number_of_increments)
previous = int(math.pow(size_i - 2, 2))
next_i = int(math.pow(size_i, 2))
numbers = next_i - previous
print("size", int(numbers / 4))
print(previous, next_i)
middle = previous + int(numbers / 8)
for i in range(previous, next_i + 1):
    if i >= middle:
        middle = middle + (int(numbers) / 4)
    if i >= input:
        print("Input", input)
        print("middle", middle)
        print("incs", number_of_increments)
        prev_middle = middle - (int(numbers) / 4)
        dist_to_middle = math.fabs(input - middle) if math.fabs(input - middle) < math.fabs(
            input - prev_middle) else math.fabs(input - prev_middle)
        print("Steps:", dist_to_middle + number_of_increments)
        break

list_of_numbers = []
index = 0
