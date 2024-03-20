<template>

    <div class='row' id='component-grid'>

          <div class='col-md-6' id='left-screen'>
            <div class='col drop-area' id='drop_0_0' :draggable='getDraggable' @dragstart="dragComponent" @drop='dropComponent' @dragover.prevent @dragenter.prevent><set-parameters id="set-parameters" :sparams="getSParams" :isDisabled="true" :calibrationState="getCalibrationState"/></div>
            <div v-if="getSyncPorts" class='col drop-area' id='drop_1_0' :draggable='getDraggable' @dragstart="dragComponent" @drop='dropComponent' @dragover.prevent @dragenter.prevent><drag-and-drop-calibrate-sync-ports id='drag-and-drop-calibrate' :calibrationState="getCalibrationState"/></div>
            <div v-else class='col drop-area' id='drop_1_0' :draggable='getDraggable' @dragstart="dragComponent" @drop='dropComponent' @dragover.prevent @dragenter.prevent><drag-and-drop-calibrate id='drag-and-drop-calibrate' :calibrationState="getCalibrationState"/></div>
          </div>

          <div class='col-md-6' id='right-screen'>
            <div class='col drop-area' id='drop_0_1' :draggable='getDraggable' @dragstart="dragComponent" @drop='dropComponent' @dragover.prevent @dragenter.prevent><data-stream id='data-stream' /></div>
          </div>

      </div>


</template>

<script>
import { mapGetters } from 'vuex';
import { dragComponent, dropComponent } from '@/js/dragAndDrop.js';

import DragAndDropCalibrate from '@/components/DragAndDropCalibrate.vue';
import DragAndDropCalibrateSyncPorts from '@/components/DragAndDropCalibrateSyncPorts.vue';
import DataStream from '../components/DataStream.vue';
import SetParameters from '../components/SetParameters.vue';

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
        SetParameters
        
    },
    computed:{
        ...mapGetters([
            'getDraggable',
            'getSParams',
            'getSyncPorts',
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