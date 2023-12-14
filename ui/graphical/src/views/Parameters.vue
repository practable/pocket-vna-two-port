<template>

    <div class='row' id='component-grid'>

          <div class='col-md-6' id='left-screen'>
            <div class='col drop-area' id='drop_0_0' :draggable='getDraggable' @dragstart="dragComponent" @drop='dropComponent' @dragover.prevent @dragenter.prevent><set-parameters id="set-parameters" :sparams="getSParams" :isDisabled="false" :calibrationState="getCalibrationState"/></div>
          </div>

          <div class='col-md-6' id='right-screen'>
            <div class='col drop-area' id='drop_0_1' :draggable='getDraggable' @dragstart="dragComponent" @drop='dropComponent' @dragover.prevent @dragenter.prevent><data-stream id='data-stream' /></div>
            
          </div>

      </div>


</template>

<script>
import { mapGetters } from 'vuex';
import { dragComponent, dropComponent } from '@/js/dragAndDrop.js';
import SetParameters from '../components/SetParameters.vue';
import DataStream from '../components/DataStream.vue';


export default {
    name: 'Parameters',
    props: [],
    data () {
        return {
            
        }
    },
    components: {
        SetParameters,
        DataStream
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
        dragComponent(event){
            dragComponent(event);
        },
        dropComponent(event){
            dropComponent(event);
        }

    }
}
</script>

<style>

</style>