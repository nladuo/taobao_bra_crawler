<template>
    <div id="app">
        <div id="main">
            <!--<basic></basic>-->
            <!--<color></color>-->
            <details></details>
        </div>
    </div>
</template>
<script type="text/ecmascript-6">
    $(".linear").height(($(window).height()));


//    import Basic from './components/BraTypeBasic.vue'
    import Details from './components/BraTypeDetails.vue'
//    import Color from './components/BraColor.vue'
    export default {
        ready(){
            var that = this;
            $.getJSON('/static/bra.json', function(data) {
                var data_basic=[],type_basic=[],data_color=[],type_color=[],data_size=[],type_size=[];
                console.log(data);
                //颜色
                $.each(data.color, function (key, word) {
                    if (key != "whole") {
                        data_color.push({"value": word, "name": key});
                        type_color.push(key);
                    }
                });
                //尺寸
                $.each(data.basic, function (key, word) {
                    if (key != "whole") {
                        data_basic.push({"value": word, "name": key + "杯"});
                        type_basic.push(key + "杯");
                    }
                });
                //百分比
                $.each(data.detail, function (key, word) {
                    var whole = data.detail.whole;
                    if (key != "whole" && word >= 1000) {
                        data_size.push((100 * word / whole).toFixed(3));
                        type_size.push(key);
                    }
                });
                console.log(data_size);
                console.log(type_size);
                that.$broadcast('data-size-loaded', data_size, type_size);
            });
        },
        components:{
//            basic: Basic,
//            color: Color,
            details: Details
        }
    }
</script>
