$(function () {
    const data = $('#bootstrap-data').data('bootstrap');
    const chart = nv.models.lineChart()
        .margin({left: 100, right: 100})
        .useInteractiveGuideline(true)
        .duration(250)
        .showLegend(true);
    const timeFormat = d3.time.format('%d %b, %I:%M %p');
    chart.xAxis
        .axisLabel('Time')
        .tickFormat(function (x) {
            return timeFormat(new Date(x * 1000));
        });
    chart.yAxis
        .axisLabel('Trump Mentions');
    d3.select('#chart-container svg')
        .datum(data['MinorMatches'])
        .call(chart);
});
