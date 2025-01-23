document.getElementById('register-form').addEventListener('submit', function(event) {
    event.preventDefault();

    var title = document.getElementById('title').value;
    var contents = document.getElementById('contents').value;
    var category = document.getElementById('category').value;
    var important = document.getElementById('important').checked;

    var data = {
        title: title,
        contents: contents,
        category: category,
        important: important
    };

    fetch('/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    .then(() => {
        window.location.href = '/'; // Redirect to home page after registration
    })
    .catch((error) => {
        window.location.href = '/'; // Redirect to home page on error
    });
});
