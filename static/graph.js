$(function () {
    const bootstrap = $('#bootstrap-data');
    const data = bootstrap.data('bootstrap');
    const chart = nv.models.lineChart()
        .margin({left: 100, right: 100, bottom: 70})
        .useInteractiveGuideline(true)
        .duration(250)
        .showLegend(true);
    const timeFormat = d3.time.format('%d %b, %I:%M %p');
    chart.xAxis
        .axisLabel('Time')
        .rotateLabels(20)
        .tickFormat(function (x) {
            return timeFormat(new Date(x * 1000));
        });
    chart.yAxis
        .axisLabel(bootstrap.data('axis'));
    d3.select('#chart-container svg')
        .datum(data[bootstrap.data('metric')])
        .call(chart);
});
