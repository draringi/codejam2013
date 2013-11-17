function updateChart (data) {
	var list = [];
	var records = data.Records;
	var len = records.length;
	for (var i = 0; i < len; i++) {
		list.push([records[i].Date, records[i].Power])
	}
	$.jqplot('futureplot', [list], {
		title:'Predicted Power Usage for the next 24 Hours',
		axes:{xaxis:{renderer:$.jqplot.DateAxisRenderer}},
		series:[{lineWidth:4, markerOptions:{show:false}}]
	});
}

function getData () {
	$.getJSON("/data", updateChart);
}

getData();
var timerID = setInterval(getData, 3*60*1000)
