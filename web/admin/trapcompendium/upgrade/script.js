var trapComponentData;

$.getJSON('/api/tc/components/short', (data) => {
  trapComponentData = data;
})

function handleTypeInput() {
  const types = [,'trigger','target','effect','universal'];
  var target = $('#component_id')[0];
  var create = document.createElement;

  return function(e) {
    $('#component_id')[0].innerHTML = '';

    if (e.target.value !== '4') {
      var type = types[e.target.value]
      var options = trapComponentData.filter(component => component.type === type);
      target.innerHTML = `
        <select name="component_id"></select>
      `;
      var menu = target.firstElementChild;

      var menuItem = create('option');
      menuItem.innerText = `All ${type}s`;
      menuItem.value = '';
      menu.append(menuItem);
      
      for (let i = 0; i < options.length; i++) {
        var menuItem = create('option');
        menuItem.innerText = options[i].name;
        menuItem.value = options[i].id;
        menu.append(menuItem);
      }
    }
  }
}

$('#type')[0].oninput = handleTypeInput();