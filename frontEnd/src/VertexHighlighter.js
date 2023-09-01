import './App.css';
import * as THREE from "three";
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js';
import { RoundedBoxGeometry } from 'three/examples/jsm/geometries/RoundedBoxGeometry.js';
import { useEffect, useState } from 'react';

const faceOffsetMap = {"F": 0, "L": 18, "B": 90, "R": 0, "U": 36, "D": 54}

const cubeSpacing = 1.1;

// eslint-disable-next-line
const colourMapping = {
  0: new THREE.Color(0xffffff),
  1: new THREE.Color(0x009b48),
  2: new THREE.Color(0xb71234),
  3: new THREE.Color(0x0046ad),
  4: new THREE.Color(0xff5800),
  5: new THREE.Color(0xffd500)
}

const defaultSideColour = new THREE.Color(0x666666);

// eslint-disable-next-line
const subCubeData = [
  {pos: new THREE.Vector3(-cubeSpacing, cubeSpacing, -cubeSpacing), faces: [[0, "U"], [9, "L"], [18, "B"]]},
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

// const createSubCubeArray = (cubeLayout) => {
//   const cubes = [];
//   subCubeData.forEach(({pos, faces}) => {
//     const cubeParts = createSubCubeParts();
//     faces.forEach(([i, f]) => {
//       colourCubeFace(f, colourMapping[cubeLayout[i]], cubeParts);
//     });
//     const cube = formSubCube(cubeParts);
//     cube.position.copy(pos);
//     cubes.push(cube);
//   });
//   return cubes;
// }

// eslint-disable-next-line
const createSubCube = () => {
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
  return {cube: cube, colours: colours};
}

// eslint-disable-next-line
const colourCubeFace = (face, colour, {cube, colours}) => {
  const o = faceOffsetMap[face];
  for (let i = o; i < o + 6 * 3; i += 3) {
    colours[i] = colour.r;
    colours[i+1] = colour.g;
    colours[i+2] = colour.b;
  }
  cube.geometry.setAttribute('color', new THREE.Float32BufferAttribute(colours, 3));
}

const App = () => {
  // eslint-disable-next-line
  const [cubeLayout, setCubeLayout] = useState([0,0,0,0,0,0,0,0,0,1,1,1,2,2,2,3,3,3,4,4,4,1,1,1,2,2,2,3,3,3,4,4,4,1,1,1,2,2,2,3,3,3,4,4,4,5,5,5,5,5,5,5,5,5]);

  // const scene = new THREE.Scene();
  // const camera = new THREE.PerspectiveCamera(75, window.innerWidth/window.innerHeight, 0.1, 1000);
  // const renderer = new THREE.WebGLRenderer();
  // renderer.setSize(window.innerWidth, window.innerHeight);

  // const cube = createSubCube();
  // colourCubeFace("F", new THREE.Color(0xaa00ff), cube);
  // // colourCubeFace("L", new THREE.Color(0xaa0000), cubeParts);
  // // colourCubeFace("B", new THREE.Color(0xffff00), cubeParts);
  // // colourCubeFace("R", new THREE.Color(0xf0f00f), cubeParts);

  // scene.add(cube.cube);

  // // createSubCubeArray(cubeLayout).forEach(c => scene.add(c));

  // const controls = new OrbitControls(camera, renderer.domElement);
  // controls.target.set(0, 0, 0);
  // controls.zoomSpeed = 0.3;
  // controls.enablePan = false;
	// controls.enableDamping = true;
	// controls.update();

  // camera.position.z = 5;

  // const animate = () => {
  //   requestAnimationFrame(animate);
  //   controls.update();
  //   renderer.render(scene, camera);
  // };
  // animate();

  

var scene = new THREE.Scene();
var raycaster = new THREE.Raycaster();

//create some camera
let camera = new THREE.PerspectiveCamera(55, window.innerWidth / window.innerHeight, 0.1, 1000);
camera.position.z = 3;
camera.lookAt(0, 0, 0);

var renderer = new THREE.WebGLRenderer({
  antialias: true
});

var controls = new OrbitControls(camera, renderer.domElement);
controls.zoomSpeed = 0.2;

renderer.setSize(window.innerWidth, window.innerHeight);
renderer.setClearColor(new THREE.Color(0x595959));

// white spotlight shining from the side, casting a shadow
var spotLight = new THREE.SpotLight(0xffffff, 2.5, 25, Math.PI / 6);
spotLight.position.set(4, 10, 7);
scene.add(spotLight);

// collect objects for raycasting, 
// for better performance don't raytrace all scene
var tooltipEnabledObjects = [];

var dodecahedronGeometry = new THREE.WireframeGeometry(new THREE.BoxGeometry().toNonIndexed());

// var dodecahedron = new THREE.Mesh(dodecahedronGeometry);
// scene.add(dodecahedron);

const wireframe = new THREE.WireframeGeometry( dodecahedronGeometry );

const line = new THREE.LineSegments( wireframe );
line.material.depthTest = false;
line.material.opacity = 0.25;
line.material.transparent = true;

scene.add( line );

var size = 0.01;
var vertGeometry = new THREE.BoxGeometry(size, size, size);
var vertMaterial = new THREE.MeshBasicMaterial({
  color: 0x0000ff,
  transparent: false
});

var verts = dodecahedronGeometry.attributes.position.array;
for (let k=0; k<verts.length; k+=3) {
  var vertMarker = new THREE.Mesh(vertGeometry, vertMaterial);

  // this is how tooltip text is defined for each box
  let tooltipText = `idx: ${k}, pos: [${verts[k].toFixed(3)},${verts[k+1].toFixed(3)},${verts[k+2].toFixed(3)}]`;
  vertMarker.userData.tooltipText = tooltipText;

  vertMarker.position.copy(new THREE.Vector3(verts[k],verts[k+1],verts[k+2]));
  scene.add(vertMarker);
  tooltipEnabledObjects.push(vertMarker);
}

function animate() {
  requestAnimationFrame(animate);
  controls.update();
  renderer.render(scene, camera);
};

// this will be 2D coordinates of the current mouse position, [0,0] is middle of the screen.
var mouse = new THREE.Vector2();

var latestMouseProjection; // this is the latest projection of the mouse on object (i.e. intersection with ray)
var hoveredObj; // this objects is hovered at the moment

// tooltip will not appear immediately. If object was hovered shortly,
// - the timer will be canceled and tooltip will not appear at all.
var tooltipDisplayTimeout;

// This will move tooltip to the current mouse position and show it by timer.
function showTooltip() {
  var divElement = document.getElementById("tooltip");
  if (divElement && latestMouseProjection) {
    divElement.style.display = "block";
    divElement.style.opacity = 1;

    var canvasHalfWidth = renderer.domElement.offsetWidth / 2;
    var canvasHalfHeight = renderer.domElement.offsetHeight / 2;

    var tooltipPosition = latestMouseProjection.clone().project(camera);
    tooltipPosition.x = (tooltipPosition.x * canvasHalfWidth) + canvasHalfWidth + renderer.domElement.offsetLeft;
    tooltipPosition.y = -(tooltipPosition.y * canvasHalfHeight) + canvasHalfHeight + renderer.domElement.offsetTop;

    var tootipWidth = divElement.offsetWidth;
    var tootipHeight = divElement.offsetHeight;

    divElement.style.left = `${tooltipPosition.x - tootipWidth/2}px`;
    divElement.style.top = `${tooltipPosition.y - tootipHeight - 5}px`;

    // var position = new THREE.Vector3();
    // var quaternion = new THREE.Quaternion();
    // var scale = new THREE.Vector3();
    // hoveredObj.matrix.decompose(position, quaternion, scale);

    const hoveredPos = hoveredObj.position;
    let points = [];
    for (let k=0; k<verts.length; k+=3) {    
      if (hoveredPos.x === verts[k] && hoveredPos.y === verts[k+1] && hoveredPos.z === verts[k+2]) {
        points.push(k);
      }
    }
    console.log(points);
    var blob = new Blob(["[" + String(points) + "]"], {type: 'text/plain'});
    var item = new window.ClipboardItem({'text/plain': blob});
    navigator.clipboard.write([item]);
    divElement.innerText = (hoveredObj.userData.tooltipText);

    setTimeout(function() {
      divElement.style.opacity = 1;
    }, 25);
  }
}

// This will immediately hide tooltip.
function hideTooltip() {
  var divElement = document.getElementById("tooltip");
  if (divElement) {
    divElement.style.display = "none";
  }
}

// Following two functions will convert mouse coordinates
// from screen to three.js system (where [0,0] is in the middle of the screen)
function updateMouseCoords(event, coordsObj) {
  coordsObj.x = ((event.clientX - renderer.domElement.offsetLeft + 0.5) / window.innerWidth) * 2 - 1;
  coordsObj.y = -((event.clientY - renderer.domElement.offsetTop + 0.5) / window.innerHeight) * 2 + 1;
}

function handleManipulationUpdate() {
  raycaster.setFromCamera(mouse, camera);
  var intersects = raycaster.intersectObjects(tooltipEnabledObjects);
  if (intersects.length > 0) {
    latestMouseProjection = intersects[0].point;
    hoveredObj = intersects[0].object;
  }

  if (tooltipDisplayTimeout || !latestMouseProjection) {
    clearTimeout(tooltipDisplayTimeout);
    tooltipDisplayTimeout = undefined;
    hideTooltip();
  }

  if (!tooltipDisplayTimeout && latestMouseProjection) {
    tooltipDisplayTimeout = setTimeout(function() {
      tooltipDisplayTimeout = undefined;
      showTooltip();
    }, 10);
  }
}

function onMouseMove(event) {
  updateMouseCoords(event, mouse);

  latestMouseProjection = undefined;
  hoveredObj = undefined;
  handleManipulationUpdate();
}

window.addEventListener('mousemove', onMouseMove, false);

animate();



  useEffect(() => {
    document.getElementById("sceneContainer").replaceChildren(renderer.domElement);
  }, [renderer.domElement]);

  return (
    <div>
      {/* <button onClick={() => {
        faceOffsetMap["F"] = faceOffsetMap["F"] + 3;
        console.log(faceOffsetMap["F"]);
        colourCubeFace("F", new THREE.Color(0xaa00ff), cube);
      }}>Hi</button> */}
      <div id="sceneContainer" />
      <div id="tooltip" />
    </div>
  );
}

export default App;
