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
        row.appendChild(createElement('td', record.patientId));
        row.appendChild(createElement('td', record.category));
        row.appendChild(createElement('td', record.path));
        table.appendChild(row);
    });
}

function createElement(tag, value) {
    const column = document.createElement(tag);
    column.appendChild(document.createTextNode(value));
    return column
}

async function startImport() {
    await fetch('./import', {method: 'PUT'});
}

async function loadCounts(elementId, countType) {
    const response = await fetch(countType + '/counts');
    const data = await response.json();

    const container = document.getElementById(elementId);
    container.appendChild(createElement('p', 'Befunde: ' + data.records));
    container.appendChild(createElement('p', 'Patienten' + data.patients));
}

window.onload = function () {
    loadCounts('database-counts', 'database');
    loadCounts('database-counts', 'filesystem');
};
