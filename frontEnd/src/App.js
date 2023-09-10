import './App.css';
import * as THREE from "three";
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js';
import { RoundedBoxGeometry } from 'three/examples/jsm/geometries/RoundedBoxGeometry.js';
import {useCallback, useEffect, useState} from 'react';

const faceOffsetMap = {"F": 672 * 3, "L": 222 * 3, "B": 822 * 3, "R": 72 * 3, "U": 372 * 3, "D": 522 * 3}

const cubeSpacing = 1;

const transforms = ["F", "L", "B", "R", "U", "D", "f", "l", "b", "r", "u", "d"];

const startingCube = [0,0,0,0,0,0,0,0,0,1,1,1,2,2,2,3,3,3,4,4,4,1,1,1,2,2,2,3,3,3,4,4,4,1,1,1,2,2,2,3,3,3,4,4,4,5,5,5,5,5,5,5,5,5];

const defaultSideColour = new THREE.Color(0x666666);

const colourMapping = {
  0: new THREE.Color(0xffffff),
  1: new THREE.Color(0x009b48),
  2: new THREE.Color(0xb71234),
  3: new THREE.Color(0x0046ad),
  4: new THREE.Color(0xff5800),
  5: new THREE.Color(0xffd500)
}

const subCubeData = [
  {pos: new THREE.Vector3(-cubeSpacing, cubeSpacing, -cubeSpacing), faces: [[0, "U"], [9, "L"], [20, "B"]]},
  {pos: new THREE.Vector3(0, cubeSpacing, -cubeSpacing), faces: [[1, "U"], [19, "B"]]},
  {pos: new THREE.Vector3(cubeSpacing, cubeSpacing, -cubeSpacing), faces: [[2, "U"], [17, "R"], [18, "B"]]},
  {pos: new THREE.Vector3(-cubeSpacing, cubeSpacing, 0), faces: [[3, "U"], [10, "L"]]},
  {pos: new THREE.Vector3(0, cubeSpacing, 0), faces: [[4, "U"]]},
  {pos: new THREE.Vector3(cubeSpacing, cubeSpacing, 0), faces: [[5, "U"], [16, "R"]]},
  {pos: new THREE.Vector3(-cubeSpacing, cubeSpacing, cubeSpacing), faces: [[6, "U"], [11, "L"], [12, "F"]]},
  {pos: new THREE.Vector3(0, cubeSpacing, cubeSpacing), faces: [[7, "U"], [13, "F"]]},
  {pos: new THREE.Vector3(cubeSpacing, cubeSpacing, cubeSpacing), faces: [[8, "U"], [14, "F"], [15, "R"]]},
  {pos: new THREE.Vector3(-cubeSpacing, 0, -cubeSpacing), faces: [[21, "L"], [32, "B"]]},
  {pos: new THREE.Vector3(-cubeSpacing, 0, 0), faces: [[22, "L"]]},
  {pos: new THREE.Vector3(-cubeSpacing, 0, cubeSpacing), faces: [[23, "L"], [24, "F"]]},
  {pos: new THREE.Vector3(0, 0, cubeSpacing), faces: [[25, "F"]]},
  {pos: new THREE.Vector3(cubeSpacing, 0, cubeSpacing), faces: [[26, "F"], [27, "R"]]},
  {pos: new THREE.Vector3(cubeSpacing, 0, 0), faces: [[28, "R"]]},
  {pos: new THREE.Vector3(cubeSpacing, 0, -cubeSpacing), faces: [[29, "R"], [30, "B"]]},
  {pos: new THREE.Vector3(0, 0, -cubeSpacing), faces: [[31, "B"]]},
  {pos: new THREE.Vector3(-cubeSpacing, -cubeSpacing, -cubeSpacing), faces: [[33, "L"], [44, "B"], [51, "D"]]},
  {pos: new THREE.Vector3(-cubeSpacing, -cubeSpacing, 0), faces: [[34, "L"], [48, "D"]]},
  {pos: new THREE.Vector3(-cubeSpacing, -cubeSpacing, cubeSpacing), faces: [[35, "L"], [36, "F"], [45, "D"]]},
  {pos: new THREE.Vector3(0, -cubeSpacing, cubeSpacing), faces: [[37, "F"], [46, "D"]]},
  {pos: new THREE.Vector3(cubeSpacing, -cubeSpacing, cubeSpacing), faces: [[38, "F"], [39, "R"], [47, "D"]]},
  {pos: new THREE.Vector3(cubeSpacing, -cubeSpacing, 0), faces: [[40, "R"], [50, "D"]]},
  {pos: new THREE.Vector3(cubeSpacing, -cubeSpacing, -cubeSpacing), faces: [[41, "R"], [42, "B"], [53, "D"]]},
  {pos: new THREE.Vector3(0, -cubeSpacing, -cubeSpacing), faces: [[43, "B"], [52, "D"]]},
  {pos: new THREE.Vector3(0, -cubeSpacing, 0), faces: [[49, "D"]]}
]

const createSubCubeArray = (cubeLayout) => {
  const cubes = [];
  subCubeData.forEach(({pos, faces}) => {
    const cube = createSubCube(pos, faces);
    faces.forEach(([i, f]) => {
      colourCubeFace(f, colourMapping[cubeLayout[i]], cube);
    });
    cube.cube.position.copy(pos);
    cubes.push(cube);
  });
  return cubes;
}

const createSubCube = (pos, faces) => {
  const geometry = new RoundedBoxGeometry();
  const positionAttribute = geometry.getAttribute( 'position' );
  const colours = [];
  for (let i = 0; i < positionAttribute.count; i += 1) {
    colours.push(defaultSideColour.r, defaultSideColour.g, defaultSideColour.b);
  }
  geometry.setAttribute('color', new THREE.Float32BufferAttribute(colours, 3));
  const material = new THREE.MeshBasicMaterial({vertexColors: true});
  material.needsUpdate = true;
  const cube = new THREE.Mesh(geometry, material);
  cube.castShadow = true;
  cube.receiveShadow = true;
  return {cube: cube, colours: colours, pos: pos, faces: faces};
}

const colourCubeFace = (face, colour, {cube, colours}) => {
  if (face === "") { //todo streamline
    for (let i = 0; i + 2 < colours.length ; i += 3) {
      colours[i] = defaultSideColour.r;
      colours[i+1] = defaultSideColour.g;
      colours[i+2] = defaultSideColour.b;
    }
  } else {
    const o =  faceOffsetMap[face];
    for (let i = o; i < o + 3 * 6; i += 3) {
      colours[i] = colour.r;
      colours[i+1] = colour.g;
      colours[i+2] = colour.b;
    }
  }
  cube.geometry.setAttribute('color', new THREE.Float32BufferAttribute(colours, 3));
}

const updateFaceColours = (cubeLayout) => {
  subCubeArray.forEach((cube) => {
    cube.faces.forEach(([i, f]) => {
      const colour = colourMapping[cubeLayout[i]];
      colourCubeFace(f, colour, cube);
    });
  });
}

const subCubeArray = createSubCubeArray(startingCube);
const scene = new THREE.Scene();
const camera = new THREE.PerspectiveCamera(75, window.innerWidth/window.innerHeight, 0.1, 1000);
camera.position.z = 5;
const renderer = new THREE.WebGLRenderer();
const controls = new OrbitControls(camera, renderer.domElement);
controls.target.set(0, 0, 0);
controls.zoomSpeed = 0.3;
controls.enablePan = false;
controls.enableDamping = true;
controls.update();

const animate = () => {
  requestAnimationFrame(animate);
  controls.update();
  renderer.render(scene, camera);
};

const App = () => {
  const [cubeLayout, setCubeLayout] = useState(startingCube);
  const [transformQueue, setTransformQueue] = useState("")
  const [playTransforms, setPlayTransforms] = useState(false)

  const transformCube = useCallback(t => {
    fetch(window.location.href + "cube", {
      method: "POST",
      body: JSON.stringify({CubeLayout: cubeLayout, Transformation: t})
    })
    .then(response => response.json())
    .then(data => {
      setCubeLayout(data);
    });
  }, [cubeLayout, setCubeLayout]);

  renderer.setSize(window.innerWidth, window.innerHeight);

  useEffect(() => {
    updateFaceColours(cubeLayout);
  }, [cubeLayout]);

  useEffect(() => {
    if (playTransforms) {
      if (transformQueue === "") {
        setPlayTransforms(false);
        return;
      }
      const timerID = setTimeout(() => {
        const transformStep = transformQueue.charAt(0);
        transformCube(transformStep);
        setTransformQueue(prev => prev.slice(1));
      }, 1000)
      return () => {
        clearTimeout(timerID)
      };
    }
  }, [playTransforms, transformQueue, transformCube])

  useEffect(() => {
    subCubeArray.forEach(c => scene.add(c.cube));
    animate();
    document.getElementById("sceneContainer").replaceChildren(renderer.domElement);
  }, []);

  return (
    <div>
      {
        transforms.map(t => <button key={t} onClick={() => {
          setPlayTransforms(false)
          transformCube(t);
        }}>{t}</button>)
      }
      <input type="text" value={transformQueue} onChange={(e) => {
        setTransformQueue(e.target.value);
        setPlayTransforms(false)
      }}/>
      <button onClick={() => setPlayTransforms(prev => !prev)} >
        {playTransforms ? "Pause": "Play"}
      </button>
      <div id="sceneContainer" />
    </div>
  );
}

export default App;
