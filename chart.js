var monAttempts = 1761;
var tueAttempts = 1714;
var wedAttempts = 1804;
var thuAttempts = 1734;
var friAttempts = 1671;
var satAttempts = 118;
var sunAttempts = 123;

window.onload = function () {

	var chart = new CanvasJS.Chart("chartContainer", {
		animationEnabled: true,
		theme: "light2", // "light1", "light2", "dark1", "dark2"
		title: {
			text: "Total Login Attempts on Days of the Week (As of 9/22/2020)"
		},
		axisX: {
			title: "Day of the Week"
		},
		axisY: {
			title: "Total Login Attempts"
		},
		data: [{
			type: "column",
			dataPoints: [
				{ y: monAttempts, label: "Monday" },
				{ y: tueAttempts, label: "Tuesday" },
				{ y: wedAttempts, label: "Wednesday" },
				{ y: thuAttempts, label: "Thursday" },
				{ y: friAttempts, label: "Friday" },
				{ y: satAttempts, label: "Saturday" },
				{ y: sunAttempts, label: "Sunday" },
			]
		}]
	});
	chart.render();

}