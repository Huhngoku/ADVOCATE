elems = 0
lines = 0
with open('dedego.log', 'r') as file:
    for line in file:
        lines += 1
        if len(line) == 0:
            continue
        elems += (line.count(';') + 1)
print(lines, elems)