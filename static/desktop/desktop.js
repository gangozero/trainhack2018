


$( document ).ready(function() {
    loadTasks();
    setInterval(function() {
        loadTasks();
    }, 5000);
});



// System section
var APIURL = ''
function getJSON(endpoint,data,func){
	var json_r = JSON.stringify(data);
    console.log("[DEBUG:JsonAPI] Request: \r\n" + json_r);
   // if(endpoint == 'tasks') return func({"tasks":[{"train":"75","carriage":"2","station":"740000210","repeat_order":false,"delivery":false,"order":[{"coffee_type":"Espresso","number":1}],"arrival_time":1540591144},{"train":"75","carriage":"2","station":"740000308","repeat_order":false,"delivery":true,"order":[{"coffee_type":"Tea","number":1},{"coffee_type":"Espresso","number":1}],"arrival_time":1540591442},{"train":"75","carriage":"2","station":"740000210","repeat_order":false,"delivery":true,"order":[{"coffee_type":"Espresso","number":2},{"coffee_type":"Cappucinno","number":1}],"arrival_time":1540591801}]});
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

function loadTasks(){
    var table = $('#table');
    getJSON('tasks',{},function(resp){
        table.html('');
        $.each(resp.tasks,function(k,task){
            var tr = $('<tr>');
            tr.append($('<td>').text(task.train));

            var dt = new Date(task.arrival_time*1000);
            tr.append($('<td>').text(dt.toTimeString().substr(0,5)));

            tr.append($('<td>').text(task.carriage));

            var delivery = 'Self pick-up'
            if(task.delivery) delivery = 'Delivery to carriage';
            tr.append($('<td>').text(delivery));

            var order = '';
            $.each(task.order,function(k,v){
                order += v.number+' x '+v.coffee_type+' <br>';
            });
            tr.append($('<td>').html(order));

            var action = "Wait"
            var diff = task.arrival_time-Math.floor(Date.now() / 1000);
            if(task.deliver && diff <= 5*60) {action = "Go to the train"; tr.addClass('go-train')}
            if(!task.deliver && diff <= 1*60) {action = "Start brewing";  tr.addClass('start-brewing')}
            if(diff < 0) {action = "Completed";  tr.removeClass('start-brewing').removeClass('go-train')}

            tr.append($('<td>').text(action));

            table.append(tr);
        });
    });

}