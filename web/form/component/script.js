function handleTypeChange(e) {
  var bool = e.target.value === "3";
  document.getElementById("param_field").setAttribute('style', `visibility: ${bool ? 'visible': 'collapse'}`);
  if (bool) {
    populateParamsTable();
    removeCostInput()
  } else {
    clearParamsTable();
    addCostInput();
  }
}

function concatParamValues(e) {
  const className = e.target.className;
  var tablerow = document.getElementsByClassName(className);
  var result = tablerow[0].value;
  for (var i = 1; i < tablerow.length; i++) {
    var value = tablerow[i].value;
    if (value !== "") {
      result += "\t" + tablerow[i].value;
    }
  }
  document.getElementsByName(className)[0].value = result;
}

function populateParamsTable() {
  const create = document.createElement;
  var table = document.getElementById('param_table')

  for (let i = 0; i < 4; i++) {

    var input = create('input');
    input.hidden = true;
    input.name = `param${i + 1}`;
    input.type = 'text';

    table.append(input);

    var row = table.insertRow(-1);

    var td = create('td');
    var paramName = create('input')
    paramName.disabled = i === 0 ? true : false;
    paramName.className = `param${i + 1}`;
    paramName.type = 'text';
    paramName.value = i === 0 ? "Cost": '';

    td.append(paramName);
    row.append(td);

    for (let j = 0; j < 7; j++) {
      td = create('td');
      var paramValue = create('input')
      paramValue.className = `param${i + 1}`;
      paramValue.type = 'text';
      paramValue.oninput = concatParamValues;
      td.append(paramValue)
      row.append(td);
    }
    table.append(row);
  }


}

function clearParamsTable() {
  var table = document.getElementById('param_table');
  table.innerHTML = '';
}

function addCostInput() {
  var costNode = document.getElementById('cost');
  if (costNode === null) {
    var target = document.getElementById('name');
    var n = document.createElement('input');
    n.required = true;
    n.type = 'number';
    n.name = n.id = 'cost';
    n.placeholder = 'Cost';
    target.insertAdjacentElement('afterend', n)
  }
}

function removeCostInput() {
  var target = document.getElementById('cost');
  target.parentNode.removeChild(target);
}

document.getElementById("type").oninput = handleTypeChange;