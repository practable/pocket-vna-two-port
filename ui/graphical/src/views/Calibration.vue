<template>

    <div class='row' id='component-grid'>

          <div class='col-md-6' id='left-screen'>
            <div v-if="getSyncPorts" class='col drop-area' id='drop_0_0' :draggable='getDraggable' @dragstart="dragComponent" @drop='dropComponent' @dragover.prevent @dragenter.prevent><drag-and-drop-calibrate-sync-ports id='drag-and-drop-calibrate' :sparams="getSParams"/></div>
            <div v-else class='col drop-area' id='drop_0_0' :draggable='getDraggable' @dragstart="dragComponent" @drop='dropComponent' @dragover.prevent @dragenter.prevent><drag-and-drop-calibrate id='drag-and-drop-calibrate' :sparams="getSParams"/></div>
          </div>

          <div class='col-md-6' id='right-screen'>
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

import DragAndDropCalibrate from '@/components/DragAndDropCalibrate.vue';
import DragAndDropCalibrateSyncPorts from '@/components/DragAndDropCalibrateSyncPorts.vue';
import DataStream from '../components/DataStream.vue';
import GraphOutput from '../components/GraphOutput.vue';
import SmithChart from '../components/SmithChart.vue';

export default {
    name: 'Calibration',
    props: [],
    data () {
        return {
            
        }
    },
    components: {
        DragAndDropCalibrate,
        DragAndDropCalibrateSyncPorts,
        DataStream,
        GraphOutput,
        SmithChart
        
    },
    computed:{
        ...mapGetters([
            'getDraggable',
            'getSParams',
            'getSyncPorts'
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