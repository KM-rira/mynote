function toggleDeleteButton() {
    var confirmDelete = document.getElementById('confirm-delete');
    var deleteBtn = document.getElementById('delete-btn');
    if (confirmDelete.checked) {
        deleteBtn.disabled = false;
        deleteBtn.style.backgroundColor = '#f44336'; // Change to red when enabled
        deleteBtn.style.cursor = 'pointer'; // Change cursor when button is enabled
    } else {
        deleteBtn.disabled = true;
        deleteBtn.style.backgroundColor = 'grey'; // Change to grey when disabled
        deleteBtn.style.cursor = 'not-allowed'; // Change cursor when button is disabled
    }
}

function deleteNote() {
    var id = document.getElementById('id').value;
    fetch('/delete', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ id: parseInt(id, 10) }) // Ensure id is sent as an integer
    })
    .then(() => {
        window.location.href = '/'; // Redirect to home page after deletion
    })
    .catch((error) => {
        console.error('Error:', error);
        window.location.href = '/'; // Redirect to home page on error
    });
}
