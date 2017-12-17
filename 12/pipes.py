import networkx as nx
graph = nx.Graph()

with open("input", "r") as f:
    connections = list(map(str.strip, f.readlines()))
set_connections = []
for c in connections:
    node, neighbors = c.split(" <-> ")
    graph.add_edges_from((node, neighbor)
                         for neighbor in neighbors.split(', '))


print ("Part 1:", len(nx.node_connected_component(graph, '0')))
print('Part 2:', nx.number_connected_components(graph))
