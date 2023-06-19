document.getElementById('registration-form').addEventListener('submit', function(event) {
    event.preventDefault();

    const formData = {
        first_name: document.getElementById('first_name').value,
        last_name: document.getElementById('last_name').value,
        email: btoa(document.getElementById('email').value), // encode email as base64
        password: document.getElementById('password').value, // encode password as base64
        password2: document.getElementById('password2').value // encode password as base64
    };

    // Retrieve CSRF token from the cookie
    const csrfToken = getCookie('CustomAPI_csrf'); 
    fetch('/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'X-CustomAPI-CSRF-Token': csrfToken // Include the CSRF token in the request headers
        },
        body: JSON.stringify(formData),
        credentials: 'include'
    }).then(function(response) {
        if (response.status == 403){
            document.getElementById('error-msg').innerHTML = "Please refresh the page and try again.";
            document.getElementById('error-msg').style.display = 'block';
        }

        if (response.status == 429){
            document.getElementById('error-msg').innerHTML = "Too many requests. Please try again later.";
            document.getElementById('error-msg').style.display = 'block';
        }

        if (response.ok) {
            document.getElementById('registration-form').reset();
            window.location.href = "html/auth/registrationsuccess.html";
        }

        // Grab the json response body
        response.json().then(function(data) {
            console.log(data);
            if (data.error){
                document.getElementById('error-msg').innerHTML = data.message;
                document.getElementById('error-msg').style.display = 'block';
            }
        })
    }).catch(function(error) {
        console.error('Error:', error);
        document.getElementById('error-msg').innerHTML = "Please refresh the page and try again. Or try clearing your browser cache.";
        document.getElementById('error-msg').style.display = 'block';
    });
    });

    function getCookie(name) {
    const cookies = document.cookie.split(';');
    for (let i = 0; i < cookies.length; i++) {
        const cookie = cookies[i].trim();
        if (cookie.startsWith(name + '=')) {
        return cookie.substring(name.length + 1);
        }
    }
    return null;
}