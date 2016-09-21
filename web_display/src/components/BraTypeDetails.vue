<template>
    <div id="details"></div>
</template>
<script>
    export default{
        events: {
            'data-size-loaded' (data_size, type_size) {

              var chart = echarts.init(this.$el);
              console.log(data_size);
              console.log(type_size);
              var option_details = {
                title: {
                  text: '各胸围尺寸所占人数柱状图',
                  /*subtext: '仅具参考意义',*/
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
                  y:'30px',
                  left: 'center',
                  data:['各尺寸所占人数百分比']
                },
                xAxis : [
                  {
                    type : 'category',
                    data : type_size
                  }
                ],
                yAxis : [
                  {
                    type : 'value',
                    name : '人数/%',
                    min: 0,
                    max: 20,
                    interval:5,
                    axisLabel : {
                      formatter: '{value}'
                    }
                  }
                ],
                series : [
                  {
                    name:'各尺寸所占人数百分比',
                    type:'bar',
                    data:data_size,
                    color:'#191970'
                  }
                ]
              };
              chart.setOption(option_details);

            }
        }
    }
</script>
