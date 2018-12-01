with open("known_input", "r") as f:
    instructions = list(map(str.strip,f.readlines()))
print(instructions)

def is_number(s):
    try:
        complex(s) # for int, long, float and complex
    except ValueError:
        return False

    return True
def parse (instructions):
    registers = {}
    instruction_index = 0
    sound_played = -1
    recovered_sound = -1
    while  instruction_index < len(instructions):
        split_ins = instructions[instruction_index].split(' ')
        print("si", split_ins)
        if split_ins[0] == 'snd':
            sound_played = int(split_ins[1]) if is_number(split_ins[1]) else  registers.get(split_ins[1],0)
        elif split_ins[0] == 'set':
            second_part = int(split_ins[2]) if is_number(split_ins[2]) else registers.get(split_ins[2],0)
            registers[split_ins[1]] =  second_part
        elif split_ins[0] == 'mod':
            second_part = int(split_ins[2]) if is_number(split_ins[2]) else registers.get(split_ins[2],0)
            registers[split_ins[1]] = registers.get(split_ins[1],0) % second_part if second_part != 0 else 0
        elif split_ins[0] == 'add':
            second_part = int(split_ins[2]) if is_number(split_ins[2]) else registers.get(split_ins[2],0)
            registers[split_ins[1]] = registers.get(split_ins[1],0) + second_part
        elif split_ins[0] == 'mul':
            second_part = int(split_ins[2]) if is_number(split_ins[2]) else registers.get(split_ins[2],0)
            registers[split_ins[1]] = int(registers.get(split_ins[1],0)) * second_part
        elif split_ins[0] == 'rcv':
            rcv = int(split_ins[1]) if is_number(split_ins[1]) else registers.get(split_ins[1],0)
            if rcv != 0:
                return sound_played
            else:
                print("no sound played")
        elif split_ins[0] == 'jgz':
            jump = ( int(split_ins[1]) if is_number(split_ins[1]) else registers.get(split_ins[1],0)) > 0
            if jump:
                instruction_index += int(split_ins[2])
                continue
            else:
                print("no jump")
#        print("index",instruction_index)
#        print("regs", registers)
        instruction_index += 1

rcv = parse(instructions)
print("sound", rcv)
