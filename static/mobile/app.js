


$( document ).ready(function() {

    showScreen('loading');
    loadMap();
});

// System section
var APIURL = ''
function getJSON(endpoint,data,func){
	var json_r = JSON.stringify(data);
    console.log("[DEBUG:JsonAPI] Request: \r\n" + json_r);

    $.post(APIURL+'/'+endpoint,json_r)
        .done(function(resp){
            console.log("[DEBUG:JsonAPI] Answer: \r\n" + resp);
            try {
                json = JSON.parse(resp);
            } catch(error) {
                console.log( "[JsonAPI] Parse Failed: \r\n" + error );
                console.log( "[JsonAPI] Recieved text: \r\n" + resp );
                json = {error:'Parse failed'}
            }
            func(json);
        })
        .fail(function( jqxhr, textStatus, error ) {
            console.log( "[JsonAPI] Request Failed: \r\n" + error );
        });
}


function showScreen(id){
    $(".screen").hide(0);
    $('#screen_'+id).fadeIn(500);
}

function getUserLocatiion(){
    if ("geolocation" in navigator){ //check geolocation available 
        //try to get user current location using getCurrentPosition() method
        navigator.geolocation.getCurrentPosition(function(position){ 
                theMap.setView([position.coords.latitude,position.coords.longitude],15)
            });
    }else{
        console.log("Browser doesn't support geolocation!");
        return [0,0]
    }
}

function getTrainID(){
    // will be code to get Train ID from beacon in the carriage
    return "75"
}

function getCarriage(){
    // will be code to get carriage number from the beacon in the carriage
    return "2"
}


// Start app section

function loadMap(){
    var map = initMap();
    getListOfStations(function(list){

        $.each(list.stations,function(k,st){
            var marker = L.marker([st.geo.lat, st.geo.lon]).addTo(map);
        });

        var i = 0;
        setStation(map,list.stations,i);
        $('#prev-station').click(function(){
            if(i > 0) i--;
            setStation(map,list.stations,i);
        });

        $('#next-station').click(function(){
            if(i < (list.stations.length-1)) i++;
            setStation(map,list.stations,i);
        });

        $('#order-coffee').click(function(){
            $('#order-buttons').fadeOut(function(){
                $('#choose-coffee').fadeIn();
                $('#add-coffee').click(function(){
                    var row = $('#coffee-type-template').clone().addClass('coffee-type');
                    $('#coffee-list').append(row);
                }).click();
                $('#go-deliver-choose').click(function(){
                    $('#choose-coffee').fadeOut(function(){
                        $('#choose-delivery').fadeIn();

                        $('#go-checkout').click(function(){
                            loadCheckout();
                        });

                        $('.deliver-picker').click(function(){
                            if($(this).hasClass('btn-primary')) return;
                            $('.deliver-picker').removeClass('btn-primary');
                            var parent = $(this).parent();
                            if(parent.data('type') == 'myself'){
                                parent.data('type','deliver');
                                $('#deliver-deliver').addClass('btn-primary'); 
                            }else{
                                parent.data('type','myself');
                                $('#deliver-myself').addClass('btn-primary'); 
                            }
                        }).click();
                    });
                });
            });
        });

        showScreen('map');
    });

    //  getUserLocatiion();
}

function setStation(map,list,num){
    var st = list[num]

    $('.station_title').text('Next station: '+st.title);
    $('.station_time').text('in ' + Math.round(st.time/60) + ' min.');
    $('.order-button').prop('disabled', true);
    $('#station-info').data('station-id',st.id).data('station-title',st.title).data('station-time',st.time);
    $('#prev-station,#next-station').prop('disabled', true);
    if(st.is_coffee == 1) $('#order-coffee').prop('disabled',false)
    if(st.is_food == 1) $('#order-food').prop('disabled',false)
    $('#order-notify').prop('disabled',false)

    if(num > 0) $('#prev-station').prop('disabled',false);
    if(num < (list.length-1)) $('#next-station').prop('disabled',false);


    map.setView([st.geo.lat,st.geo.lon],15);
}


function initMap(){
	var mapCont = $('#map');

	var theMap = L.map(mapCont.attr('id'), {attributionControl: false});
	L.control.attribution({prefix:''}).addAttribution('&copy; Gangozero | &copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>').addTo(theMap);  

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
		maxZoom: 19
	}).addTo(theMap);
    theMap.setView([0,0], 19);
    return theMap;
}

function getListOfStations(callback){
    var train = getTrainID();
    var carriage = getCarriage();
    getJSON('list',{train:train,carriage:carriage},function(resp){
        callback(resp)
    })
}

function loadCheckout(){
    showScreen('loading');
    var train = getTrainID();
    var carriage = getCarriage();
    var station_id = $('#station-info').data('station-id');
    var repeat = $('#repeat').is(':checked');
    var deliver = false;
    if($('#deliver-type').data('type') == 'deliver') deliver= true;
    var order = {}
    $('.coffee-type').each(function(){
        var item = $(this).val();
        if(item in order)
            order[item] = order[item]+1;
        else
            order[item] = 1;

    });
    var coffeelist = []
    $.each(order,function(type,number){
        coffeelist.push({coffee_type:type,number:number})
    })
    getJSON('order',{train:train,carriage:carriage,station:station_id,repeat_order:repeat,delivery:deliver,caffeeshop_id:1,order:coffeelist},function(resp){
        showScreen('payment');
        setTimeout(function(){
            $('#qrcode').qrcode("order_id:"+Math.floor(Math.random() * 100) + 1  );
            showScreen('success');
        }, 2000);
    })
}


