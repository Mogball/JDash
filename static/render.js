dat.GUI.prototype.removeFolder = function (name) {
    var folder = this.__folders[name];
    if (!folder) {
        return;
    }
    folder.close();
    this.__ul.removeChild(folder.domElement.parentNode);
    delete this.__folders[name];
};

const bootstrap = $('#bootstrap-data');

var gui = new dat.GUI();
var scene = new THREE.Scene();
var files = bootstrap.data('bootstrap');
var renderObj = {
    file: "main_battle_tank"
};
var savedState = undefined;

var camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 10, 50);
camera.position.x = 10;
camera.position.y = 10;
camera.position.z = 13;
camera.lookAt(0, 0, 0);

var renderer = new THREE.WebGLRenderer({antialias: true});
renderer.setPixelRatio(window.devicePixelRatio);
renderer.setSize(window.innerWidth, window.innerHeight);
renderer.setClearColor(0x000000, 1);
document.body.appendChild(renderer.domElement);

var ambientLight = new THREE.AmbientLight(0x000000);
scene.add(ambientLight);

var lights = [];
lights[0] = new THREE.PointLight(0xffffff, 1, 0);
lights[1] = new THREE.PointLight(0xffffff, 1, 0);

lights[0].position.set(0, 200, 0);
lights[1].position.set(100, 200, 100);

scene.add(lights[0]);
scene.add(lights[1]);

var loader = new THREE.JSONLoader();

function loadFile(file) {
    loader.load("/static/models/" + file + ".json", onLoad);
}

loadFile(renderObj.file);

var refresh = undefined;

function onLoad(geometry) {
    var mesh = new THREE.Mesh(geometry);
    mesh.material = new THREE.MeshStandardMaterial({
        color: 0x2194ce,
        flatShading: true,
        roughness: 0.72,
        metalness: 0.72
    });
    if (savedState) {
        mesh.material.color.setHex(savedState.color);
        mesh.material.roughness = savedState.roughness;
        mesh.material.metalness = savedState.metalness;
        mesh.material.wireframe = savedState.wireframe;
        mesh.material.flatShading = savedState.flatShading;
        mesh.rotation.y = savedState.rotation;
    }

    scene.add(mesh);
    guiMesh(gui, mesh);
    guiMaterial(gui, mesh.material);

    function render() {
        refresh = requestAnimationFrame(render);
        renderer.render(scene, camera);
        mesh.rotation.y += 0.005;
    }

    window.addEventListener('resize', function () {
        camera.aspect = window.innerWidth / window.innerHeight;
        camera.updateProjectionMatrix();
        renderer.setSize(window.innerWidth, window.innerHeight);
    }, false);
    render()
}

function guiMesh(gui, mesh) {
    var folder = gui.addFolder('Mesh');
    folder.add(renderObj, 'file', files).onChange(function (val) {
        if (refresh) {
            cancelAnimationFrame(refresh);
        }
        savedState = {
            color: mesh.material.color.getHex(),
            roughness: mesh.material.roughness,
            metalness: mesh.material.metalness,
            wireframe: mesh.material.wireframe,
            flatShading: mesh.material.flatShading,
            rotation: mesh.rotation.y
        };
        gui.removeFolder('Mesh');
        gui.removeFolder('Material');
        scene.remove(mesh);
        loadFile(val);
    });
    folder.open();
}

function guiMaterial(gui, material) {
    var folder = gui.addFolder('Material');
    var data = {
        color: material.color.getHex(),
        wireframe: material.wireframe
    };
    folder.addColor(data, 'color').onChange(function (val) {
        if (typeof val === 'string') {
            val = val.replace('#', '0x');
        }
        material.color.setHex(val);
    });
    var roughness = folder.add(material, 'roughness', 0.0, 1.0);
    var metalness = folder.add(material, 'metalness', 0.0, 1.0);
    folder.add(data, 'wireframe').onChange(function (val) {
        material.wireframe = val;
        material.flatShading = !val;
        material.roughness = val ? 1.0 : 0.72;
        material.metalness = val ? 0.5 : 0.72;
        material.needsUpdate = true;
        roughness.updateDisplay();
        metalness.updateDisplay();
    });
}
