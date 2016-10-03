<style>
    #detail {
        width: 70%;
        height: 530px;
        margin: 0 auto;
    }
</style>
<template>
    <div id="detail"></div>
</template>
<script>
    export default{
        props:{
            bra: Object
        },
        events: {
            'reload-detail-chart' () {
                $("#detail").width(document.body.clientWidth*0.7);
//                console.log('reload-detail-chart');
                var chart = echarts.init(this.$el);
                var option_detail = {
                    title: {
                        text: '各胸围尺寸所占人数柱状图',
                        x:'center'
                    },
                    tooltip: {
                        trigger: 'axis'
                    },
                    toolbox: {
                        show: true,
                        feature: {
                            mark: {show: true},
                            dataView: {show: true, readOnly: false},
                            magicType: {show: true, type: ['line', 'bar']},
                            restore: {show: true},
                            saveAsImage : {show: true}
                        }
                    },
                    calculable : true,
                    legend: {
                        orient: 'horizontal',
                        y: '30px',
                        left: 'center',
                        data:['各尺寸所占人数百分比']
                    },
                    xAxis : [
                        {
                            type : 'category',
                            data : this.bra.type_detail
                        }
                    ],
                    yAxis : [
                        {
                            type: 'value',
                            name: '人数/%',
                            min: 0,
                            max: 20,
                            interval:5,
                            axisLabel: {
                                formatter: '{value}'
                            }
                        }
                    ],
                    series : [
                        {
                            name: '各尺寸所占人数百分比',
                            type: 'bar',
                            data: this.bra.detail,
                            color: '#191970'
                        }
                    ]
                };
                chart.setOption(option_detail);
            }
        }
    }
</script>
