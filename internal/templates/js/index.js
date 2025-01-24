document.addEventListener('DOMContentLoaded', () => {
    // JSONデータを取得
    const notesJSON = document.getElementById('notes-data').textContent;
    const notes = JSON.parse(notesJSON); // JSON文字列をオブジェクトに変換
    const tbody = document.getElementById('notes-body');

    if (notes.length === 0) {
        const noDataRow = document.createElement('tr');
        const noDataCell = document.createElement('td');
        noDataCell.colSpan = 8;
        noDataCell.textContent = "No notes found.";
        noDataRow.appendChild(noDataCell);
        tbody.appendChild(noDataRow);
    } else {
        notes.forEach(note => {
            const row = document.createElement('tr');

            // Selectボタン
            const selectCell = document.createElement('td');
            const selectButton = document.createElement('button');
            selectButton.textContent = "Select";
            selectButton.onclick = () => alert(`Selected Note ID: ${note.ID}`);
            selectCell.appendChild(selectButton);
            row.appendChild(selectCell);

            // ID
            const idCell = document.createElement('td');
            idCell.textContent = note.ID;
            row.appendChild(idCell);

            // Title
            const titleCell = document.createElement('td');
            titleCell.textContent = note.title;
            row.appendChild(titleCell);

            // Contents
            const contentsCell = document.createElement('td');
            contentsCell.textContent = note.contents;
            row.appendChild(contentsCell);

            // Category
            const categoryCell = document.createElement('td');
            categoryCell.textContent = note.category;
            row.appendChild(categoryCell);

            // Important
            const importantCell = document.createElement('td');
            const importantCheckbox = document.createElement('input');
            importantCheckbox.type = "checkbox";
            importantCheckbox.disabled = true;
            if (note.important) {
                importantCheckbox.checked = true;
            }
            importantCell.appendChild(importantCheckbox);
            row.appendChild(importantCell);

            // Created At
            const createdAtCell = document.createElement('td');
            createdAtCell.textContent = new Date(note.createdAt).toLocaleString();
            row.appendChild(createdAtCell);

            // Updated At
            const updatedAtCell = document.createElement('td');
            updatedAtCell.textContent = new Date(note.updatedAt).toLocaleString();
            row.appendChild(updatedAtCell);

            tbody.appendChild(row);
        });
    }
});
