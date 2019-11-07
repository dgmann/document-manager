async function validate() {
    const response = await fetch('./validate');
    const data = await response.json();
    drawValidationErrorList(data.errors);
    drawRecordTable(data.resolvables);
}

function drawValidationErrorList(errors) {
    const list = document.getElementById("validation-errors");
    errors.forEach(error => {
        const item = document.createElement('li');
        item.appendChild(document.createTextNode(error));
        list.appendChild(item);
    });
}

function drawRecordTable(records) {
    const table = document.getElementById("record-table");
    records.forEach(record => {
        const row = document.createElement('tr');
        row.appendChild(createColumn(record.patientId));
        row.appendChild(createColumn(record.category));
        row.appendChild(createColumn(record.path));
        table.appendChild(row);
    });
}

function createColumn(value) {
    const column = document.createElement('td');
    column.appendChild(document.createTextNode(value));
    return column
}
