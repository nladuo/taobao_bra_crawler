<template>
    <div id="app">
        <div id="main">
            <basic :bra.sync="bra" v-bind:style="{ display: now_page == 0 ? 'block' : 'none'}"></basic>
            <color :bra.sync="bra" v-bind:style="{ display: now_page == 1 ? 'block' : 'none'}"></color>
            <detail :bra.sync="bra" v-bind:style="{ display: now_page == 2 ? 'block' : 'none'}"></detail>
        </div>
        <nav>
            <ul class="pager">
                <li class="previous" @click="previous">
                    <a href="#">
                        <span class="glyphicon glyphicon-arrow-left" aria-hidden="true" style="font-size: 22px"></span>
                    </a>
                </li>
                <li class="next" @click="next">
                    <a href="#">
                        <span class="glyphicon glyphicon-arrow-right" aria-hidden="true" style="font-size: 22px"></span>
                    </a>
                </li>
            </ul>
        </nav>
    </div>
</template>
<style>
    #app, #main{
        height: 100%;
    }
</style>
<script type="text/ecmascript-6">
    $(".linear").height(($(window).height()));


    import Basic from './components/BraBasic.vue'
    import Detail from './components/BraDetail.vue'
    import Color from './components/BraColor.vue'
    export default {
        data(){
            return{
                now_page: 0,
                bra:{
                    color:[],
                    type_color:[],
                    basic:[],
                    type_basic:[],
                    detail:[],
                    type_detail:[]
                }
            }
        },
        methods:{
            reload_chart: function () {
                switch (this.now_page){
                    case 0: this.$broadcast('reload-basic-chart');break;
                    case 1: this.$broadcast('reload-color-chart');break;
                    case 2: this.$broadcast('reload-detail-chart');break;
                }
            },
            next: function () {
                this.now_page = (this.now_page + 1) % 3;
                this.reload_chart();
            },
            previous: function () {
                this.now_page = (this.now_page + 2) % 3;
                this.reload_chart();
            }
        },
        ready(){
            var that = this;
            $.getJSON('static/bra.json', function(data) {
//                console.log(data);
                var bra = that.bra;
                //按罩杯分类
                $.each(data.basic, function (key, word) {
                    if (key != "whole") {
                        bra.basic.push({"value": word, "name": key + "杯"});
                        bra.type_basic.push(key + "杯");
                    }
                });
                //按颜色分类
                $.each(data.color, function (key, word) {
                    if (key != "whole") {
                        bra.color.push({"value": word, "name": key});
                        bra.type_color.push(key);
                    }
                });
                //按具体的罩杯分类
                $.each(data.detail, function (key, word) {
                    var whole = data.detail.whole;
                    if (key != "whole" && word >= 1000) {
                        bra.detail.push((100 * word / whole).toFixed(3));
                        bra.type_detail.push(key);
                    }
                });
                that.bra = bra;
                that.reload_chart();

            });
        },
        components:{
            basic: Basic,
            color: Color,
            detail: Detail
        }
    }
</script>
