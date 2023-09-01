import re

def getFace(index):
    if index <= 8:
        return "U"
    if index >= 45:
        return "D"
    if index in [9,10,11,21,22,23,33,34,35]:
        return "L"
    if index in [12,13,14,24,25,26,36,37,38]:
        return "F"
    if index in [15,16,17,27,28,29,39,40,41]:
        return "R"
    if index in [18,19,20,30,31,32,42,43,44]:
        return "B"
    print("Invalid index: " + str(index))

def getFacePairs(indexes):
    indexes = map(int, indexes.split(", "))
    indexes = map(lambda i: "[" + str(i) + ", \"" + getFace(i) + "\"]" , indexes)
    return ", ".join(indexes)

sheets_data = "(-1, 1, -1)|[0, 9, 20]	(0, 1, -1)|[1, 19]	(1, 1, -1)|[2, 17, 18]	(-1, 1, 0)|[3, 10]	(0, 1, 0)|[4]	(1, 1, 0)|[5, 16]	(-1, 1, 1)|[6, 11, 12]	(0, 1, 1)|[7, 13]	(1, 1, 1)|[8, 14, 15]	(-1, 0, -1)|[21, 32]	(-1, 0, 0)|[22]	(-1, 0, 1)|[23, 24]	(0, 0, 1)|[25]	(1, 0, 1)|[26, 27]	(1, 0, 0)|[28]	(1, 0, -1)|[29, 30]	(0, 0, -1)|[31]	(-1, -1, -1)|[33, 44, 51]	(-1, -1, 0)|[34, 48]	(-1, -1, 1)|[35, 36, 45]	(0, -1, 1)|[37, 46]	(1, -1, 1)|[38, 39, 47]	(1, -1, 0)|[40, 50]	(1, -1, -1)|[41, 42, 53]	(0, -1, -1)|[43, 52]	(0, -1, 0)|[49]"

subCubes = []

for subCube in re.findall(r"\(.*?]", sheets_data):
    position, indexes = subCube.split("|")
    position = position.replace("1", "cubeSpacing")
    subCubeString = "{pos: new THREE.Vector3" + position
    subCubeString += ", faces: [" + getFacePairs(indexes.replace("[", "").replace("]", "")) + "]}"
    subCubes.append(subCubeString)

print("[\n  " + ",\n  ".join(subCubes) + "\n]")
