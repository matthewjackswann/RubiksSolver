MAX_MOVES = 20  # god's number for the 3x3


def factorial(n):
    return 1 if n == 0 else n * factorial(n - 1)


TOTAL_PERMUTATIONS = factorial(8) * pow(3, 7) * (factorial(12) // 2) * pow(2, 11)  # according to wikipedia

node_n = 0
graph = []  # node -> [nodes]

with open("graph.csv", "r") as graph_file:
    node_no = 0
    for line in graph_file:
        edges = line.strip().split(',')
        if node_n == 0:
            node_n = len(edges)
        else:
            assert node_n == len(edges)

        node_edges = []
        for node_to, exists in enumerate(edges):
            if exists != '_':
                node_edges.append(node_to)
        graph.append(node_edges)
        node_no += 1

cube_n = 1  # count solved cube

node_count = [0 for _ in range(node_n)]
# assumes node 0 is init node
node_count[0] = 1

for m in range(MAX_MOVES):
    next_node_count = [0 for _ in range(node_n)]

    for node in range(node_n):
        for next_node in graph[node]:
            # adds a new cube for cube that can move to this node in 1 turn
            next_node_count[next_node] += node_count[node]

    node_count = next_node_count
    cube_n += sum(node_count)

print(f'Cubes generated: {float(cube_n)}')
print(f'Permutations covered: {float(cube_n * 6)}')
print(f'Total permutations: {float(TOTAL_PERMUTATIONS)}')
print(f'Excess permutations (multiplier): {(cube_n * 6) / TOTAL_PERMUTATIONS}')
