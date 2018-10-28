<template>
  <div id="app">
    <div id="main">
      <bra-basic :bra.sync="bra" v-bind:style="{ display: now_page == 0 ? 'block' : 'none'}"></bra-basic>
      <bra-color :bra.sync="bra" v-bind:style="{ display: now_page == 1 ? 'block' : 'none'}"></bra-color>
      <bra-detail :bra.sync="bra" v-bind:style="{ display: now_page == 2 ? 'block' : 'none'}"></bra-detail>
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
<script type="text/ecmascript-6">
  import BraBasic from './components/BraBasic.vue'
  import BraDetail from './components/BraDetail.vue'
  import BraColor from './components/BraColor.vue'
  export default {
    components:{
      BraBasic, BraDetail, BraColor
    },
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
      reload_chart() {
        switch (this.now_page){
          case 0: this.$broadcast('reload-basic-chart');break;
          case 1: this.$broadcast('reload-color-chart');break;
          case 2: this.$broadcast('reload-detail-chart');break;
        }
      },
      next() {
        this.now_page = (this.now_page + 1) % 3;
        this.reload_chart();
      },
      previous() {
        this.now_page = (this.now_page + 2) % 3;
        this.reload_chart();
      }
    },
    ready(){
      $(".linear").height(($(window).height()));
      let that = this;
      $.getJSON('static/bra.json', function(data) {
        var bra = that.bra;
        //按罩杯分类
        $.each(data.basic, function (key, word) {
          if (key !== "whole") {
            bra.basic.push({"value": word, "name": key + "杯"});
            bra.type_basic.push(key + "杯");
          }
        });
        //按颜色分类
        $.each(data.color, function (key, word) {
          if (key !== "whole") {
            bra.color.push({"value": word, "name": key});
            bra.type_color.push(key);
          }
        });
        //按具体的罩杯分类
        $.each(data.detail, function (key, word) {
          let whole = data.detail.whole;
          if (key !== "whole" && word >= 1000) {
            bra.detail.push((100 * word / whole).toFixed(3));
            bra.type_detail.push(key);
          }
        });
        that.bra = bra;
        that.reload_chart();
      });
    }
  }
</script>
<style scoped>
  @import './assets/style.css';
  #app, #main{
    height: 100%;
  }
</style>
