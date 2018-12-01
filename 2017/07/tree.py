import re
import json

with open("input_known", "r") as f:
    content = f.read()
    regex = re.compile(
        "((?P<id>\w+)\s\((?P<weight>\d+)\))(?:$| -> )(?P<children>.*)", re.MULTILINE)
    programs = {}
    for prog in [m.groupdict() for m in regex.finditer(content)]:
        programs[prog['id']] = prog
        children = prog.get('children').split(",")
        programs[prog['id']]['children'] = [x.strip()
                                            for x in children if len(x) > 0]
    for id, prog in programs.items():
        for c in prog.get('children'):
            programs[c]['parent'] = prog['id']
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
    print(parent_program.get('id'))
