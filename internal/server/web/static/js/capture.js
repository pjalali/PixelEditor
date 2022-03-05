// Based on code from Mozilla: https://developer.mozilla.org/en-US/docs/Web/API/WebRTC_API/Taking_still_photos

let resolutions = [
    // {Width: 320, Height: 240},
    {Width: 640, Height: 480},
    // {Width: 1024, Height: 768},
    {Width: 1280, Height: 720},
    {Width: 1920, Height: 1080}
];


function testWebcam() {
    resolutions.forEach(resolution => {
        navigator.mediaDevices.getUserMedia({video: { width: resolution.Width, height: resolution.Height  }, audio: false})
            .then(function(stream) {
                let {width, height} = stream.getTracks()[0].getSettings();
                console.log(`${width}x${height}`); // 640x480
            })
            .catch(function(err) {
                console.log("An error occurred while starting webcam: " + err);
            });
    });
}

// Citation: https://developer.mozilla.org/en-US/docs/Web/API/WebRTC_API/Taking_still_photos
function startup() {
    testWebcam();
    var streaming = false;
    var video = document.getElementById('video');
    var canvas = document.getElementById('canvas');
    var startbutton = document.getElementById('startbutton');

    navigator.mediaDevices.getUserMedia({video: { width: {ideal: 640}, height: {ideal: 360}  }, audio: false})
        .then(function(stream) {
            video.srcObject = stream;
            video.play();
        })
        .catch(function(err) {
            console.log("An error occurred while starting webcam: " + err);
            alert("An error occurred while starting webcam. Please try again.");
        });

    video.addEventListener('canplay', function() {
        if (!streaming) {
            let width = video.videoWidth;
            let height = video.videoHeight;
        
            video.setAttribute('width', width);
            video.setAttribute('height', height);
            canvas.setAttribute('width', width);
            canvas.setAttribute('height', height);
            document.getElementById("wData").value = width;
            document.getElementById("hData").value = height;

            streaming = true;
        }
    }, false);

    startbutton.addEventListener('click', function(ev) {
        takepicture();
        ev.preventDefault();
    }, false);
}


// https://developer.mozilla.org/en-US/docs/Web/API/WebRTC_API/Taking_still_photos
function takepicture() {
    var canvas = document.getElementById('canvas');
    var context = canvas.getContext('2d');

    context.drawImage(video, 0, 0, canvas.width, canvas.height);

    // https://stackoverflow.com/a/24289420
    var base64ImgData = canvas.toDataURL('image/png');

    document.getElementById("imgData").value = base64ImgData;

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