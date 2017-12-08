import re
import json
from copy import deepcopy
import operator


with open("input", "r") as f:
    content = f.read()
    regex = re.compile(
        "((?P<id>\w+)\s\((?P<weight>\d+)\))(?:$| -> )(?P<children>.*)", re.MULTILINE)
    programs = {}
    hits = [m.groupdict() for m in regex.finditer(content)]
    for prog in hits:
        programs[prog['id']] = prog
        children = prog.get('children').split(",")
        programs[prog['id']]['children'] = [x.strip()
                                            for x in children if len(x) > 0]
    for id, prog in programs.items():
        for c in prog.get('children'):
            programs[c]['parent'] = prog['id']
    all_clear = False
    saved_programs = deepcopy(programs)
    while True:
        to_delete = []
        for id, prog in programs.items():
            if len(prog.get('children')) == 0:
                to_delete.append(prog.get('id'))
            else:
                new_children = []
                for child in prog.get('children'):
                    if child not in to_delete and child in programs:
                        new_children = new_children + [child]
                prog['children'] = new_children
        for p in to_delete:
            del programs[p]
        if len(programs) == 1:
            parent_program = programs.get(list(programs.keys())[0])
            break
    saved_programs[parent_program.get(
        'id')]['parent'] = parent_program.get('id')

    def mark_unbalanced(node, programs):
        w_node = int(node.get('weight'))
        if node.get('children'):
            children_weight = []
            for c in node.get('children'):
                w = mark_unbalanced(programs[c], programs)
                children_weight.append(w)
            if len(set(children_weight)) > 1:
                node["children_balanced"] = False
                for c in node.get('children'):
                    nc = saved_programs[c]
            else:
                node["children_balanced"] = True
            node["cw"] = sum(children_weight)
            sw_child = node["cw"]
        else:
            sw_child = 0
        return w_node + sw_child
    root = saved_programs[parent_program.get('id')]
    mark_unbalanced(root, saved_programs)
    discrepancy_found = False

    def is_discrepant(child_node, parent_node):
        return child_node.get('children_balanced', True) != parent_node.get('children_balanced', True) and parent_node.get('id') != parent_node.get('parent')

    def find_discrepancy(node, saved_programs):
        for c in node.get("children"):
            child = saved_programs[c]
            if is_discrepant(child, node):
                weights = {}
                for c in node.get("children"):
                    choild = saved_programs[c]
                    weight = choild['cw'] + int(choild['weight'])
                    weights[weight] = weights.get(weight, 0) + 1
                inv_weights = {v: k for k, v in weights.items()}
                off_value = inv_weights[1]
                good_value = inv_weights[
                    [x for x in inv_weights.keys() if x != 1][0]]
                for c in node.get("children"):
                    choild = saved_programs[c]
                    weight = choild['cw'] + int(choild['weight'])
                    if weight == off_value:
                        print("New weight",
                              int(choild['weight']) + good_value - off_value)
                return node
            for cl in child.get('children', []):
                a = find_discrepancy(saved_programs[cl], saved_programs)
                if a is not None:
                    return a
        return None

    print ("discrepant", find_discrepancy(root, saved_programs))
