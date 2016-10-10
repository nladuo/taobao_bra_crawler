<style>
  #basic {
    width: 70%;
    height: 530px;
    margin: 0 auto;
  }
</style>
<template>
  <div id="basic"></div>
</template>
<script>
  export default{
    props:{
      bra: Object
    },
    events: {
      'reload-basic-chart' () {
        $("#basic").width(document.body.clientWidth*0.7);
        var chart = echarts.init(document.getElementById('basic'));
        var option_basic = {
          title : {
            text: '淘宝网用户购买文胸尺寸饼状图',
            subtext: '仅具参考意义',
            x:'center'
          },
          tooltip : {
            trigger: 'item',
            formatter: "{a} <br/>{b} : {c} ({d}%)"
          },
          legend: {
            orient: 'horizontal',
            y:'60px',
            left: 'center',
            data: this.bra.type_basic
          },
          series : [
            {
              name: '该尺寸购买人数',
              type: 'pie',
              radius : '55%',
              center: ['50%', '60%'],
              data: this.bra.basic,
              color:['#EEB422','#708090','#FFA07A','#8B658B','#EE3B3B','#6B8E23','#5F9EA0','#FFC125','#696969'],
              itemStyle: {
                emphasis: {
                  shadowBlur: 10,
                  shadowOffsetX: 0,
                  shadowColor: 'rgba(0, 0, 0, 0.5)'
                }
              }
            }
          ]
        };
        chart.setOption(option_basic);
      }
    }
  }
</script>
