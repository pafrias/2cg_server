/**
 * Created by journ on 11/21/2017.
 */

var SVM = {viewModel: new SpellViewModel()};

$( document ).ready(function() {
    ko.applyBindings(SVM.viewModel, document.getElementById('spell-list-wrapper'));
});

function SpellModel(data) {
    //console.log(data);
    var self = this;
    self.spellId = data.spell_id;
    self.spellName = data.spell_name;
    self.spellClass = ko.computed(function (d) {
        return arrayToString(data.classes);
    });
    self.spellTier = data.tiers;
    self.spellDetails = data.description;
    self.spellSchool = data.schools;
    self.spellCastTime = data.cast_times;
    self.spellRange = data.ranges;
    self.spellComponents = ko.computed(function (d) {
        return arrayToString(data.components);

    });
    self.spellDurations = data.duration;
    self.spellTags = ko.computed(function (d) {
        return arrayToString(data.tags);
    });
    self.spellCon = ko.computed(function (d) {
        return (data.requires_concentration === true) ? 'Yes' : 'No' ;
    });

    self.spellMechanics = ko.computed(function (d) {
        return arrayToString(
            [
                (data.requires_verbal === true) ? 'V' : '',
                (data.requires_somatic === true) ? 'S' : '',
                (data.requires_material === true) ? 'M' : ''
            ]
        );
    });
}

function SpellViewModel(){
    var self = this;
    self.spellList = ko.observableArray([]);
    self.noSearchResults = ko.observable(false);
    self.emptyList = ko.observable(true);

    self.getSpellList = function (uri) {
        $.getJSON(uri, function (data) {
            var mappedObjects = $.map(data, function (item) {
                return new SpellModel(item);
            });
            self.spellList(mappedObjects);
        })
            .always(function (data) {
                if(self.spellList.length !== 0) self.emptyList(false);
            });
    };

    self.updateList = function (data) {
        removeCarets(['name', 'tier', 'class']);
        var mappedItems = $.map(data, function (item) {
            return new SpellModel(item);
        });
        self.spellList(mappedItems);
        if(self.spellList().length > 0 ){
            self.noSearchResults(false);
            self.emptyList(false);
        }else{
            self.emptyList(false);
            self.noSearchResults(true);
        }
    };

    self.onCollapse = function (item, e) {
        //console.log(e);
        var bodyId = '#' + e.target.id;
        //console.log(bodyId);
        var body = $(bodyId);
        body.find('.card').toggleClass('border-dark border-success');
        body.toggleClass('mb-1');
        var headerId = body.data('header');
        //console.log(headerId);
        var header = $(headerId);
        header.toggleClass('border-bottom-0');
        header.addClass('mb-1');
        header.toggleClass('border-dark border-success');
        header.removeClass('selected');
        header.find('.fa').toggleClass('fa-plus fa-minus');
        return true;
    };

    self.onExpand = function (item, e) {
        //console.log(e);
        var bodyId = '#' + e.target.id;
        //console.log(bodyId);
        var body = $(bodyId);
        body.find('.card').toggleClass('border-dark border-success');
        body.toggleClass('mb-1');
        var headerId = body.data('header');
        //console.log(headerId);
        var header = $(headerId);
        header.toggleClass('border-bottom-0');
        header.toggleClass('border-dark border-success');
        header.addClass('selected');
        header.removeClass('mb-1');
        header.find('.fa').toggleClass('fa-plus fa-minus');
        return true;
    };

    self.sortByName = function (item, e) {

        var a = $('#accordion-header-name').find('.fa');

        if(a.hasClass('ascending')){
            a.toggleClass('ascending descending');
            self.spellList.sort(function (l,r) {
                return l.spellName < r.spellName;
            });
            a.removeClass('fa-caret-up').addClass('fa-caret-down');
        }else{
            a.toggleClass('ascending descending');
            self.spellList.sort(function (l,r) {
                return l.spellName > r.spellName;
            });
            a.removeClass('fa-caret-down').addClass('fa-caret-up');
        }

        removeCarets(['class', 'tier']);
    };

    self.sortByClass = function () {

        var a = $('#accordion-header-class').find('.fa');

        if(a.hasClass('ascending')){
            a.toggleClass('ascending descending');
            self.spellList.sort(function (l,r) {
                return l.spellClass() < r.spellClass();
            });
            a.removeClass('fa-caret-up').addClass('fa-caret-down');

        }else{
            a.toggleClass('ascending descending');
            self.spellList.sort(function (l,r) {
                return l.spellClass() > r.spellClass();
            });
            a.removeClass('fa-caret-down').addClass('fa-caret-up');

        }

        removeCarets(['name', 'tier']);
    };

    self.sortByTier = function () {
        var a = $('#accordion-header-tier').find('.fa');

        if(a.hasClass('ascending')){
            a.removeClass('ascending').addClass('descending');
            self.spellList.sort(function (l,r) {
                return l.spellTier < r.spellTier;
            });
            a.removeClass('fa-caret-up').addClass('fa-caret-down');
        }else{
            a.removeClass('descending').addClass('ascending');
            self.spellList.sort(function (l,r) {
                return l.spellTier > r.spellTier;
            });

            a.removeClass('fa-caret-down').addClass('fa-caret-up');
        }

        removeCarets(['name', 'class']);
    };

    self.onLoad = function () {
        self.getSpellList("./cache/response.json");
    };

    self.onLoad();
}


function removeCarets(array) {
    var headers = [];
    array.forEach(function (value, index) {
        headers.push($('#accordion-header-' + value));
    });

    headers.forEach(function (value, index) {
        value.find('.fa').removeClass('fa-caret-up').removeClass('fa-caret-down');
    });
}

function arrayToString(array) {
    if(array != null){
        var worker = $.grep(array, function (val) {
            return val != '';
        });
        return worker.join(", ");
    }
    return '';
}

function getSearchParams() {
    var filters = [];
    $('.selected-param').each(function () {
        var table = $(this).data('type');
        var a = {
            table: table,
            parameter: $(this).data('value')
        };
        filters.push( a );
    });

    //console.log(filters);
    if(filters.length < 1){
        return '';
    }else{
        return JSON.stringify(filters);
    }
}

function searchApi() {
    var params = getSearchParams();
    var uri = './cache/response.json';
    if(params.length){
        $.getJSON(uri, {search_parameters:getSearchParams()}, function (data) {
            //console.log(data);
            SVM.viewModel.updateList(data);
        })
        .fail(function (data) {
            error.log(data);
        });
    }else{
        $.getJSON(uri,function (data) {
            //console.log(data);
            SVM.viewModel.updateList(data);
        })
        .fail(function (data) {
            error.log(data);
        });
    }
}