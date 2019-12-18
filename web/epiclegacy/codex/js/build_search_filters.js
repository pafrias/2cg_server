$(document).ready(function (e) {
    ko.applyBindings(new SpellComponentViewModel(), document.getElementById('search-filter-wrapper'));
});

function SpellComponentModel(data) {
    var self = this;
    self.value = data.value;
    self.type = data.type;
    self.constrained = data.constrained;
}

function SpellComponentViewModel() {
    var self = this;
    self.filterList = ko.observableArray([]);
    self.selectedParameters = ko.observableArray([]);
    self.classList = ko.observableArray([]);
    self.tierList = ko.observableArray([]);
    self.schoolList = ko.observableArray([]);
    self.rangeList = ko.observableArray([]);
    self.durationList = ko.observableArray([]);
    self.castTimeList = ko.observableArray([]);
    self.hasParameters = ko.observable(false);

    self.getList = function (uri, list, type, constrained) {
        $.getJSON(uri, function (data) {
            var mappedData = $.map(data, function (item) {
                return new SpellComponentModel({value:item, type:type, constrained:constrained});
            });
            list(mappedData);
        });
    };

    self.getList("./cache/classes.json", self.classList, 'class', false);
    self.getList("./cache/tiers.json", self.tierList, 'tier', true);
    self.getList("./cache/schools.json", self.schoolList, 'school', true);
    self.getList("./cache/ranges.json", self.rangeList, 'range', true);
    self.getList("./cache/durations.json", self.durationList, 'duration', true);
    self.getList("./cache/cast_times.json", self.castTimeList, 'cast_time', true);


    self.filterSelected = function (obj, evt) {
        var tid = '#' + evt.target.id;
        var element = $(tid);
        var type = element.data('filter');
        switch (type){
            case 'class':
                self.filterList(self.classList());
                break;
            case 'tier':
                self.filterList(self.tierList());
                break;
            case 'school':
                self.filterList(self.schoolList());
                break;
            case 'range':
                self.filterList(self.rangeList());
                break;
            case 'duration':
                self.filterList(self.durationList());
                break;
            case 'cast':
                self.filterList(self.castTimeList());
                break;
            case 'mechanics':
                self.filterList(new SpellComponentModel({value:"Concentration", type:'concentration', constrained:false}));
                break;
            default:
                self.filterList([]);
        }
        return true;
    };

    self.addKeywordParameter = function (obj, evt) {
        var element = $('#search-input');
        var value = element.val();
        if(value !== ''){
            element.val('');
            var o = {
                value:value,
                type:'keyword',
                constrained:false
            };
            self.selectedParameters.push(new SpellComponentModel(o));
            self.hasParameters(true);
            searchApi();
        }
    };

    self.addParameter = function (obj, evt) {
        if(!containsValue(self.selectedParameters(), obj)){
            if(!containsType(self.selectedParameters(), obj)){
                self.selectedParameters.push(new SpellComponentModel({value:obj.value, type:obj.type, constrained:obj.constrained}));
                self.hasParameters(true);
                searchApi();
            }
        }
        var tid = "button[id='" + evt.target.id + "']";
        var element = $(tid);
        element.trigger('activeChanged');
        return true;
    };

    self.removeParameter = function (obj, evt)  {
        self.selectedParameters.remove(obj);
        if(self.selectedParameters().length === 0){
            self.hasParameters(false);
        }
        searchApi();
        return true;
    };

    self.removeActive = function (obj, evt) {
        $('.btn-filter').removeClass('active');
    };

}

function containsValue(array, obj){
    var flag = false;
    $.each(array, function (key, item) {
        if(obj.value === item.value){
            flag = true;
            flashExistingParameter(item.value);
        }
    });
    return flag;
}

function containsType(array, obj) {
    var flag = false;
    if(obj.constrained === true){
        $.each(array, function (key, item) {
            if(obj.type === item.type){
                flag = true;
                flashExistingParameter(item.value);
            }
        });
    }
    return flag;
}

function flashExistingParameter(value){
    var id = "button[id='search-param-" + value +"']";
    var element = $(id);
    element.toggleClass('btn-outline-success btn-outline-danger');
    element.fadeTo(300, 0.3, function() {
        $(this).fadeTo(500, 1.0).delay(100).toggleClass('btn-outline-success btn-outline-danger');
    });
}