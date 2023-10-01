import matplotlib.pyplot as plt
import networkx as nx
from matplotlib.lines import Line2D

# plt.style.use('dark_background')

with open("graph.csv", "r") as f:
    lines = list(map(lambda x: x.strip().split(","), f.readlines()))

G = nx.DiGraph()
G.add_nodes_from(range(len(lines)))
node_colours = ["r"]
node_colours.extend(["b" for _ in range(len(lines) - 1)])

pos = nx.circular_layout(G)

transform_colour_map = {
    'F': "r",
    'U': "b",
    'B': "r",
    'D': "c",
    'L': "m",
    'R': "y"
}
edge_colours = []
edge_styles = []

for i in range(len(lines)):
    for j in range(len(lines[i])):
        # i -> j
        if lines[i][j] == '_':
            continue
        G.add_edge(i, j)
        transform = lines[i][j]
        if transform.isupper():
            edge_styles.append("--")
        else:
            edge_styles.append("-")
        edge_colours.append(transform_colour_map[transform.upper()])

nx.draw_networkx_nodes(G, pos, node_color=node_colours)
nx.draw_networkx_edges(G, pos, edge_color=edge_colours, style=edge_styles)

legend_lines = [Line2D([0], [0], color=c, lw=2) for c in transform_colour_map.values()]
legend_titles = [t + " clockwise rotation" for t in transform_colour_map.keys()]

legend_lines.append(Line2D([0], [0], lw=0))
legend_titles.append("")

legend_lines.append(Line2D([0], [0], color='k', lw=2, linestyle="dashed"))
legend_titles.append("anticlockwise rotation")

legend_lines.append(Line2D([], [], color="white", marker='o', markersize=20, markerfacecolor="r"))
legend_titles.append("Starting State")

plt.legend(legend_lines, legend_titles, loc='upper left')
plt.show()
