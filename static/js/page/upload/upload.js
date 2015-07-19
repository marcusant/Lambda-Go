var UPLOAD_URL = '/upload';
var RESPONSE_URL_BASE = '/';
var DISPLAY_URL_BASE = 'https://lambda.sx'
var sizeLimit = 20; // MB
var apikey = '';
var allowedTypes = ['png'];

var dropZone = document.body;
var uploadTitle = document.getElementById('uploadTitle');
var selectInput = document.getElementById('chooseFile');
var finishedUploads = document.getElementById('finishedUploads');

// Show copy icon on hover
document.body.addEventListener('dragover', function(e) {
  e.stopPropagation();
  e.preventDefault();
  e.dataTransfer.dropEffect = 'copy';
})

// On file drop
document.body.addEventListener('drop', function(e) {
  e.stopPropagation();
  e.preventDefault(); // stop the browser from redirecting
  var files = e.dataTransfer.files;
  for(var i = 0; i < files.length; i++) {
    var file = files[i];
    checkAndUpload(file);
  }
})

selectInput.addEventListener('change', function(e) {
  checkAndUpload(selectInput.files[0]);
})

function onUploadFinish(responseText) {
  console.log(responseText);
  var response = JSON.parse(responseText);
  if(response.success) {
    var url = RESPONSE_URL_BASE + response.files[0].url;

    // Append finished upload entry
    var entry = document.createElement('a');
    entry.href = url;
    var urlDiv = document.createElement('div');
    urlDiv.className = "finished-uploads";
    urlDiv.innerHTML = DISPLAY_URL_BASE + url;
    entry.appendChild(urlDiv);
    finishedUploads.appendChild(entry);

    finishedUploads.hidden = false;
  } else {
    var errors = 'Upload failed due to the following errors:\n';
    for(var j = 0; j < errors.length; j++) {
      errors.append(response.errors[j] + '\n');
    }
    alert(errors);
  }
}

function checkAndUpload(file) {
  if(typeAllowed(file)) {
    if(file.size <= sizeLimit*1000000) {
      uploadFile(file, onUploadFinish);
    } else {
      alert('File is too big. Max filesize is ' + sizeLimit + ' MB.')
    }
  } else {
    alert('Filetype "' + file.type + '" is not supported!');
  }
}

function uploadFile(file, onFinish) {
  var xhr = new XMLHttpRequest();
  var fd = new FormData();
  xhr.open('POST', UPLOAD_URL, true);
  fd.append('apikey', apikey)
  fd.append('file', file);
  createStatusIndicator(xhr, file);
  xhr.onreadystatechange = function() { // on upload finish
    if(xhr.readyState == 4 && xhr.status == 200) {
      onFinish(xhr.responseText);
    }
  }
  xhr.send(fd);
}

function createStatusIndicator(xhr, file) {
  var progressSection = document.getElementById('uploadProgress');

  var circleContainer = document.createElement('div');
  circleContainer.className = 'upload-circle-container';
  var outerCircle = document.createElement('div');
  outerCircle.className = 'upload-percent-circle';
  var innerCircle;
  if(isImage(file)) {
    innerCircle = document.createElement('img');
    innerCircle.className = 'upload-top-circle-img';
    innerCircle.src = URL.createObjectURL(file);
  } else {
    innerCircle = document.createElement('div')
    innerCircle.className = 'upload-top-circle';
  }

  outerCircle.appendChild(innerCircle);
  circleContainer.appendChild(outerCircle);
  progressSection.appendChild(circleContainer);

  xhr.upload.addEventListener('progress', function(e) {
    var pct = e.loaded / e.total;
    if(pct >= 100) {
      progressSection.removeChild(circleContainer);
    } else {
      var degrees = pct*360;

      // I cannot explain for the life of me why this works, but it does, so don't touch it unless you're web jesus
      if(degrees <= 180) {
        degrees += 90;

        outerCircle.style.backgroundImage =
            'linear-gradient(' + degrees +
            'deg, transparent 50%, #F48FB1 50%),' +
            'linear-gradient(90deg, #F48FB1 50%, transparent 50%)';
      } else {
        degrees -= 90;

        outerCircle.style.backgroundImage =
            'linear-gradient(90deg, transparent 50%, #AD1457 50%),' +
            'linear-gradient(' + degrees +
            'deg, #F48FB1 50%, transparent 50%)';
      }
    }
  }, false);
}

function isImage(file) {
  return file.type.lastIndexOf('image/', 0) === 0; // beginsWith('image/')
}

function typeAllowed(file) {
  for(var i = 0; i < allowedTypes.length; i++) {
    var t = allowedTypes[i].toLowerCase();
    if(endsWith(file.name.toLowerCase(), '.' + t)) {
      return true;
    }
  }
  return false;
}

function endsWith(str, suffix) {
  return str.indexOf(suffix, str.length - suffix.length) !== -1;
}
