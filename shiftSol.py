baseString = "0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53".split(",")

compString = "0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,42,43,44,33,34,35,36,37,38,39,40,41,51,48,45,52,49,46,53,50,47".split(",")

links = {}

for i in range(len(baseString)):
    if baseString[i] != compString[i]:
        links[baseString[i]] = compString[i]


def getLinks(linkMap, root, end):
    # print(root, end)
    nextNode = linkMap[root]
    if nextNode == end:
        return []
    else:
        return [nextNode] + getLinks(linkMap, nextNode, end)


coveredNodes = []

for node in links:
    if node not in coveredNodes:
        chain = getLinks(links, node, node)
        chain.append(node)
        coveredNodes.extend(chain)
        chain = chain[::-1]
        for i in range(len(chain)):
            print(chain[i] + ": " + chain[(i + 1) % len(chain)] + ",")
