<template>

    <div class='row' id='component-grid'>

          <div class='col-md-6' id='left-screen'>
            <div class='col drop-area' id='drop_0_0' :draggable='getDraggable' @dragstart="dragComponent" @drop='dropComponent' @dragover.prevent @dragenter.prevent><set-parameters :sparams="getSParams" :isDisabled="false" :calibrationState="getCalibrationState"/></div>
          </div>

          <div class='col-md-6' id='right-screen'>
            <div class='col drop-area' id='drop_0_1' :draggable='getDraggable' @dragstart="dragComponent" @drop='dropComponent' @dragover.prevent @dragenter.prevent><data-stream id='data-stream' /></div>
            
          </div>

      </div>


</template>

<script>
import { mapGetters } from 'vuex';




export default {
    name: 'Parameters',
    props: [],
    data () {
        return {
            
        }
    },
    components: {
        
        
    },
    computed:{
        ...mapGetters([
            'getDraggable',
            'getSParams',
        ])
    },
    watch:{
        
    },
    created(){
        
    },
    mounted(){
        

    },
    methods:{
        

    }
}
</script>

<style>

</style>