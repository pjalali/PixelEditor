// Based on code from Mozilla: https://developer.mozilla.org/en-US/docs/Web/API/WebRTC_API/Taking_still_photos

var width = 320;    // We will scale the photo width to this
var height = 0;     // This will be computed based on the input stream


// Citation: https://developer.mozilla.org/en-US/docs/Web/API/WebRTC_API/Taking_still_photos
function startup() {
    var streaming = false;
    var video = document.getElementById('video');
    var canvas = document.getElementById('canvas');
    var startbutton = document.getElementById('startbutton');

    navigator.mediaDevices.getUserMedia({video: true, audio: false})
    .then(function(stream) {
        video.srcObject = stream;
        video.play();
    })
    .catch(function(err) {
        console.log("An error occurred: " + err);
    });

    video.addEventListener('canplay', function(ev){
        if (!streaming) {
        height = video.videoHeight / (video.videoWidth/width);
    
        video.setAttribute('width', width);
        video.setAttribute('height', height);
        canvas.setAttribute('width', width);
        canvas.setAttribute('height', height);
        streaming = true;
        }
    }, false);

    startbutton.addEventListener('click', function(ev){
        takepicture();
        ev.preventDefault();
    }, false);
}

function setDimensions() {
    var photo = document.getElementById('photo');
    photo.style.width = document.getElementsByName('width').value;
    photo.style.height = document.getElementsByName('height').value;
}

function resize(w,h) {
    width = w;
    startup();
    console.log(height)
    document.getElementById("video").style.width=w+"px";
    document.getElementById("video").style.height=h+"px";
}

function uploadImage() {
    var canvas = document.getElementById('canvas');
    var ctx = canvas.getContext('2d');

    fi = document.getElementById("fileInput");

    if (fi.files && fi.files[0]) {
        if (fi.files[0].type !== "image/png") {
            alert("Please only upload .png images.");
            return;
        }
        form = document.getElementById('uploadImage')
        form.submit()
    }
}

// https://developer.mozilla.org/en-US/docs/Web/API/WebRTC_API/Taking_still_photos
function takepicture() {
    var canvas = document.getElementById('canvas');
    var context = canvas.getContext('2d');

    canvas.width = width;
    canvas.height = height;
    context.drawImage(video, 0, 0, width, height);

    // https://stackoverflow.com/a/24289420
    var data = canvas.toDataURL('image/png');

    var imgData = document.getElementById("imgData");
    imgData.value = data;

    var wData = document.getElementById("wData");
    wData.value = width;

    var hData = document.getElementById("hData");
    hData.value = height;

    sizeData.value = width;

    document.getElementById("imgDataForm").submit(); 
}

function reset() {
    document.getElementById("photo").src = document.getElementById("originalImage").value;
    document.getElementById("rOffset").value = 0;
    document.getElementById("gOffset").value = 0;
    document.getElementById("bOffset").value = 0;
    document.getElementById("contrast").value = 0;
    document.getElementById("hue").value = 0;
    document.getElementById("sat").value = 0;
    document.getElementById("light").value = 0;
}