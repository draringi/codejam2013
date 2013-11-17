function updateChart (data) {
	var list = [];
	var records = data.records;
	var len = records;
	for (var i = 0; i < len; i++) {
		list.push([records[i].Date, records[i].Power])
	}
	$.jqplot('futureplot', [list], {
		title:'Predicted Power Usage for the next 24 Hours',
		axes:{xaxis:{renderer:$.jqplot.DateAxisRenderer}},
		series:[{lineWidth:4, markerOptions:{style:'square'}}]
	});
}

function getData () {
	$.getJSON("/data", updateChart);
}
