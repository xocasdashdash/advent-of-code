with open("input", "r") as f:
    registers = f.readlines()

base_dict = {
    "reg": "",
    "operation": "",
    "operand": "",
    "condition": ""
}

reg_values = {}
values = []


def apply_operation(register, operation,  operand):
    curr_value = reg_values.get(register, 0)
    if operation == "inc":
        new_value = curr_value + int(operand)
    else:
        new_value = curr_value - int(operand)
    reg_values[register] = new_value
    values.append(new_value)


def condition(condition, operation, reg,  operand):
    register_to_find = condition.split(" ")[0]
    cond = condition.split(" ")[1]
    val_to_compar = condition.split(" ")[2]
    curr_value = reg_values.get(register_to_find, 0)
    if curr_value == 0:
        reg_values[register_to_find] = 0
    if cond == "==" and int(val_to_compar) == curr_value:
        apply_operation(reg, operation, operand)
    elif cond == "!=" and int(val_to_compar) != curr_value:
        apply_operation(reg, operation, operand)
    elif cond == "<=" and curr_value <= int(val_to_compar):
        apply_operation(reg, operation, operand)
    elif cond == "<" and curr_value < int(val_to_compar):
        apply_operation(reg, operation, operand)
    elif cond == ">" and curr_value > int(val_to_compar):
        apply_operation(reg, operation, operand)
    elif cond == ">=" and curr_value >= int(val_to_compar):
        apply_operation(reg, operation, operand)


for reg in registers:
    new_reg = base_dict.copy()
    parsed_line = reg.split(" ")
    new_reg["reg"] = parsed_line[0]
    new_reg["operation"] = parsed_line[1]
    new_reg["operand"] = parsed_line[2]
    new_reg["condition"] = reg[(reg.index("if ") + 3):]. strip()
    # Find values
    condition(new_reg["condition"], new_reg["operation"],
              new_reg["reg"], new_reg["operand"])

print("mv", max(reg_values.values()))
print("mhv", max(values))
