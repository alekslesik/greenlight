document.addEventListener('DOMContentLoaded', function () {
    fetch("http://localhost:4000/v1/healthcheck").then(function (response) {
        response.text().then(function (text) {
            document.getElementById("output").innerHTML = text;
        });
    },
        function (err) {
            document.getElementById("output").innerHTML = err;
        }
    );
});