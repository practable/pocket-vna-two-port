<template>

    <div class='row' id='component-grid'>

          <div class='col-lg-6' id='left-screen'>
            <div class='col drop-area' id='drop_0_0' :draggable='getDraggable' @dragstart="dragComponent" @drop='dropComponent' @dragover.prevent @dragenter.prevent><set-parameters id="set-parameters" :sparams="getSParams" :isDisabled="true" :calibrationState="getCalibrationState"/></div>
            <div class='col drop-area' id='drop_1_0' :draggable='getDraggable' @dragstart="dragComponent" @drop='dropComponent' @dragover.prevent @dragenter.prevent><drag-and-drop-measurement id='dragAndDropMeasurement' :calibrationState='getCalibrationState' /></div>
             <div class='col drop-area' id='drop_2_0' :draggable='getDraggable' @dragstart="dragComponent" @drop='dropComponent' @dragover.prevent @dragenter.prevent><download id='download' :sparams="getSParams" /></div>
          </div>

          <div class='col-lg-6' id='right-screen'>
            <div class='col drop-area' id='drop_0_1' :draggable='getDraggable' @dragstart="dragComponent" @drop='dropComponent' @dragover.prevent @dragenter.prevent><data-stream id='data-stream' /></div>
            <div class='col drop-area' id='drop_1_1' :draggable='getDraggable' @dragstart="dragComponent" @drop='dropComponent' @dragover.prevent @dragenter.prevent><smith-chart id='smith-chart' :sparams="getSParams" /></div>
            <div class='col drop-area' id='drop_2_1' :draggable='getDraggable' @dragstart="dragComponent" @drop='dropComponent' @dragover.prevent @dragenter.prevent><graph-output id='11' type="single" sparams="s11"/></div>
            <div class='col drop-area' id='drop_3_1' :draggable='getDraggable' @dragstart="dragComponent" @drop='dropComponent' @dragover.prevent @dragenter.prevent><graph-output id='12' type="single" sparams="s12"/></div>
            <div class='col drop-area' id='drop_4_1' :draggable='getDraggable' @dragstart="dragComponent" @drop='dropComponent' @dragover.prevent @dragenter.prevent><graph-output id='21' type="single" sparams="s21"/></div>
            <div class='col drop-area' id='drop_5_1' :draggable='getDraggable' @dragstart="dragComponent" @drop='dropComponent' @dragover.prevent @dragenter.prevent><graph-output id='22' type="single" sparams="s22"/></div>
          </div>

      </div>


</template>

<script>
import { mapGetters } from 'vuex';
import { dragComponent, dropComponent } from '@/js/dragAndDrop.js';

import GraphOutput from '../components/GraphOutput.vue';
import DataStream from '../components/DataStream.vue';
import DragAndDropMeasurement from '../components/DragAndDropMeasurement.vue';
import SmithChart from '../components/SmithChart.vue';
import Download from '../components/Download.vue';
import SetParameters from '../components/SetParameters.vue';

export default {
    name: 'Measurement',
    props: [],
    data () {
        return {
            
        }
    },
    components: {
        GraphOutput,
        DataStream,
        DragAndDropMeasurement,
        SmithChart,
        Download,
        SetParameters
    },
    computed:{
        ...mapGetters([
            'getDraggable',
            'getSParams',
            'getCalibrationState'
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