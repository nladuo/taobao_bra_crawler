<template>
  <div id="color"></div>
</template>
<script>
  export default{
    props:{
      bra: Object
    },
    events: {
      'reload-color-chart' () {
        $("#color").width(document.body.clientWidth*0.7);
        let chart = echarts.init(document.getElementById('color'));
        let option_color = {
          title : {
            text: '淘宝网用户购买文胸颜色饼状图',
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
            data: this.bra.type_color
          },
          series : [
            {
              name: '该颜色购买人数',
              type: 'pie',
              radius : '55%',
              center: ['50%', '60%'],
              data: this.bra.color,
              color: ['#EEB422','#BEBEBE','#FFA07A','#8B658B','#EE3B3B','#6B8E23','#5F9EA0','#FFC125','#696969'],
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
        chart.setOption(option_color);
      }
    }
  }
</script>
<style scoped>
  #color {
    width: 70%;
    height:530px;
    margin: 0 auto;
  }
</style>
