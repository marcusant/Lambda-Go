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
  xhr.onreadystatechange = function() { // on upload finish
    if(xhr.readyState == 4 && xhr.status == 200) {
      onFinish(xhr.responseText);
    }
  }
  xhr.send(fd);
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
